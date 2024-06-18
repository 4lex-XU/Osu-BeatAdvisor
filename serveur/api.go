package main

import (
	"beatadvisor/serveur/data"
	"beatadvisor/serveur/handlers"
	"beatadvisor/serveur/structures"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
   ----------------------------------------------------
       Response/Requests Structures for the osu API
   ----------------------------------------------------
*/

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type BeatmapResponse struct {
	Beatmaps []structures.Beatmap `json:"beatmaps"`
}

// les structures des beatmaps de l'api osu
type APIBeatmap struct {
	Difficulty_rating float32       `json:"difficulty_rating"`
	Id                int           `json:"id"`
	Mode              string        `json:"mode"`
	Status            string        `json:"status"`
	Total_length      int           `json:"total_length"`
	User_id           int           `json:"user_id"`
	Url               string        `json:"url"`
	Beatmapset        APIBeatmapset `json:"beatmapset"`
}

type APIBeatmapset struct {
	Title          string `json:"title"`
	Language       string `json:"language"`
	Tags           string `json:"tags"`
	Ranked_date    string `json:"ranked_date"`
	Submitted_date string `json:"submitted_date"`
}

/*
   ------------------------
       GLOBAL VARIABLES
   ------------------------
*/

var (
	// Nom de la bd
	dbName = data.DbName

	// variables pour extraction des beatmaps
	cpt  = 100  // l'id des beatmaps à extraire, il sera incrémenté à chaque requête vers l'api
	flag = true // si on n'a récupéré aucune beatmap, ce flag sera mis à false pour pouvoir avancer cpt de 100 et éviter de boucler

	// token api pour accès aux beatmaps
	apitoken string

	// pour les beatmaps
	genres    = []string{"video game", "anime", "rock", "pop", "electronic", "metal", "classical", "jazz", "hip hop", "folk", "funk"}
	languages = []string{"japanese", "english", "chinese", "french", "german", "korean", "spanish", "italian", "russian", "instrumental"}
)

/*
   ----------------------
       MAIN FUNCTION
   ----------------------
*/

func init() {
	data.Init()
}

func router() http.Handler {
	router := http.NewServeMux()

	// authentification handlers
	router.HandleFunc("/api/user/register", handlers.HandleRegister)
	router.HandleFunc("/api/user/login", handlers.HandleLogin)
	router.HandleFunc("/api/user/logout", handlers.HandleLogout)

	// users handlers
	router.HandleFunc("/api/user/", handlers.HandleGetUser)           // must have user_id
	router.HandleFunc("/api/user/delete/", handlers.HandleDeleteUser) // must have user_id

	// playlist handlers
	router.HandleFunc("/api/playlist/create", handlers.HandleCreatePlaylist)
	router.HandleFunc("/api/playlist/", handlers.HandleGetPlaylist)                          // must have playlist_id
	router.HandleFunc("/api/playlist/delete/", handlers.HandleDeletePlaylist)                // must have playlist_id
	router.HandleFunc("/api/playlist/edit/", handlers.HandleEditPlaylist)                    // must have playlist_id
	router.HandleFunc("/api/playlist/add", handlers.HandleAddPlaylistForUser)                // playlist_id and user_id in the body
	router.HandleFunc("/api/playlist/get/", handlers.HandleGetPlaylistsForUser)              // must have user_id
	router.HandleFunc("/api/playlist/delete/beatmap/", handlers.HandleDeleteBeatmapPlaylist) // must have playlist_id, beatmap_id in the body
	router.HandleFunc("/api/playlist/get/all", handlers.HandleGetAllPlaylists)

	// social handlers
	router.HandleFunc("/api/playlist/comment/", handlers.HandleAddComment)           // must have playlist_id
	router.HandleFunc("/api/playlist/comment/delete/", handlers.HandleDeleteComment) // must have playlist_id/comment_id
	router.HandleFunc("/api/playlist/like/", handlers.HandleLike)                    // must have playlist_id
	router.HandleFunc("/api/playlist/like/delete/", handlers.HandleDeleteLike)       // must have comment_id
	router.HandleFunc("/api/playlist/share/", handlers.HandleSharePlaylist)          // must have playlist_id
	router.HandleFunc("/api/user/edit", handlers.HandleEditProfile)
	router.HandleFunc("/api/user/friends", handlers.HandleAddFriend)            // must have friend_id in the body
	router.HandleFunc("/api/user/friends/", handlers.HandleCheckFollow)         // must have user_id
	router.HandleFunc("/api/user/friends/delete/", handlers.HandleRemoveFriend) // must have friend_id

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	// Pour le lancement sur le web
	buildDir := "../client/build"

	fileServer := http.FileServer(http.Dir(buildDir))

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := buildDir + r.URL.Path

		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Si le fichier n'existe pas, servir index.html
			http.ServeFile(w, r, buildDir+"/index.html")
		} else {
			fileServer.ServeHTTP(w, r)
		}
	}))

	return c.Handler(router)
}

func main() { // compléter les handlers

	router := router()

	// fetch token
	startTokenRoutine()
	// fetch beatmaps
	startFetchRoutine()

	fmt.Println("Connecté à MongoDB !")

	//test create user
	// fmt.Println("create")
	// TestCreateUser()

	fmt.Println("Serveur démarré sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func TestCreateUser() {
	if data.Client == nil {
		log.Fatal("client MongoDB non initialisé")
	}
	collection := data.Client.Database(dbName).Collection("Users")

	newUser := bson.D{
		{Key: "email", Value: "test@example.com"},
		{Key: "pseudo", Value: "TestUser"},
		{Key: "password", Value: "testPassword"},
	}

	// insertion du nouvel utilisateur
	_, err := collection.InsertOne(context.Background(), newUser)
	if err != nil {
		fmt.Println("Erreur lors de la création de l'utilisateur:", err)
		return
	}

	fmt.Println("TestCreateUser passé avec succès")
}

/*
   ------------------------------------
       BEATMAP EXTRACT FROM OSU API
   ------------------------------------
*/

// fetchNewToken récupère un nouveau token d'authentification pour l'API osu, appeler toutes les 23h59
func fetchNewToken() {
	fmt.Println("Récupération d'un nouveau token en cours...")

	// client id secret pour obtenir un token
	data := "client_id=31082&client_secret=H5rzIds3OdHwHje1AP1ISxA3RyVSNgmIOv8F8H2b&grant_type=client_credentials&scope=public"
	s := strings.NewReader(data)

	// requête POST /oauth/token de l'api externe
	req, err := http.NewRequest("POST", "https://osu.ppy.sh/oauth/token", s)
	if err != nil {
		fmt.Println("Erreur lors de la création de la requête :", err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la requête pour obtenir le token :", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la réponse :", err)
		return
	}

	// le token pour pouvoir récupérer les beatmaps
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		fmt.Println("Erreur lors du décodage de la réponse JSON :", err)
		return
	}

	apitoken = tokenResp.AccessToken
	//fmt.Println("Token: ", apitoken)
}

func startTokenRoutine() {
	go func() {
		for {
			fetchNewToken()
			time.Sleep(23*time.Hour + 59*time.Minute)
		}
	}()
}

/* fetchBeatmaps récupère les beatmaps de l'API osu! et les insère dans la base de données,
   on peut régler la fréquence des requêtes */

// pour l'instant x beatmaps récupérées toutes les y secondes, peut être augmenter y quand on arrive aux beatmaps plus récentes (4 million-eme)
// On augmente l'ID de x à chaque requête et on commence par le dernier ID de la dernière beatmap
// il faut faire attention à ce que la base de données a au moins 1 beatmaps sinon mettre le flag a false en debut pour la première execution
func fetchBeatmaps() {
	fmt.Println("Récupération des beatmaps en cours...")
	collection := data.Client.Database(dbName).Collection("Beatmaps")
	x := 50
	var beatmapIDs []string // Liste pour stocker les IDs à récupérer

	if flag {
		// commencer par le dernier id de la dernière beatmap
		var lastBeatmap structures.Beatmap
		err := collection.FindOne(context.Background(), bson.D{}, options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})).Decode(&lastBeatmap) // le dernier de la collection
		if err != nil {
			fmt.Println("Erreur lors de la récupération du dernier ID de beatmap:", err)
			return
		}
		cpt = lastBeatmap.Id
	} else {
		cpt += x
	}

	fmt.Println("Démarrage à partir de l'ID de beatmap:", cpt)

	// Remplissage de la liste des IDs (50 par 50)
	for i := 0; i < x; i++ {
		beatmapIDs = append(beatmapIDs, strconv.Itoa(cpt))
		cpt++
	}

	// Construction de l'URL avec les IDs collectés
	requestURL := "https://osu.ppy.sh/api/v2/beatmaps?" // utiliser la requete getBeatmaps et non getBeatmap
	for _, id := range beatmapIDs {
		requestURL += "ids%5B%5D=" + id + "&"
	}
	requestURL = strings.TrimRight(requestURL, "&") // Enlever le dernier '&'

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		fmt.Println("Erreur lors de la création de la requête :", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+apitoken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Exécution de la requête
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de l'appel à l'API osu! :", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la réponse :", err)
		return
	}

	var response BeatmapResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Erreur lors du décodage de la réponse JSON :", err)
		return
	}

	var allBeatmaps []interface{}
	// json -> beatmap
	for _, apiBeatmap := range response.Beatmaps {
		collection := data.Client.Database(dbName).Collection("Beatmaps")
		count, err := collection.CountDocuments(context.Background(), bson.M{"id": apiBeatmap.Id})
		if err != nil {
			fmt.Println("La beatmap existe déjà:", err)
			continue
		}
		// Si la beatmap n'existe pas déjà dans la base de donnée
		if count == 0 {
			genre := ""
			language := ""

			// pour video game (genre composé de 2 mots)
			tag_compose := " " + strings.Join(strings.Split(apiBeatmap.Beatmapset.Tags, " "), " ") + " "
			for _, g := range genres {
				if genre == "" && strings.Contains(tag_compose, " "+g+" ") {
					genre = g
					break
				}
			}

			tags := strings.Split(apiBeatmap.Beatmapset.Tags, " ")
			for _, lang := range languages {
				if contains(tags, lang) {
					language = lang
					break
				}
			}

			newBeatmap := structures.Beatmap{
				Difficulty_rating: apiBeatmap.Difficulty_rating,
				Id:                apiBeatmap.Id,
				Mode:              apiBeatmap.Mode,
				Status:            apiBeatmap.Status,
				Total_length:      apiBeatmap.Total_length,
				User_id:           apiBeatmap.User_id,
				Url:               apiBeatmap.Url,
				Beatmapset: structures.Beatmapset{
					Title:          apiBeatmap.Beatmapset.Title,
					Language:       language,
					Genre:          genre,
					Tags:           apiBeatmap.Beatmapset.Tags,
					Ranked_date:    apiBeatmap.Beatmapset.Ranked_date,
					Submitted_date: apiBeatmap.Beatmapset.Submitted_date,
				},
			}
			allBeatmaps = append(allBeatmaps, newBeatmap)
		}
	}

	_, err = collection.InsertMany(context.Background(), allBeatmaps)
	if err == nil {
		flag = true
		fmt.Println("Beatmaps insérées avec succès")
	} else {
		flag = false
		fmt.Println("Rien à récupérer")
	}

}

func startFetchRoutine() {
	go func() {
		for {
			fetchBeatmaps()
			// sleep après chaque execution
			time.Sleep(10 * time.Minute)
		}
	}()
}

// Fonction helper pour vérifier si une liste contient une chaîne
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if strings.EqualFold(v, str) {
			return true
		}
	}
	return false
}
