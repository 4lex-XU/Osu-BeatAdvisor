package handlers

import (
	"beatadvisor/serveur/data"
	"beatadvisor/serveur/structures"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
   ------------------------------------
       Response/Requests Structures
   ------------------------------------
*/

// create playlist request
type PlaylistRequest struct {
	Genres      []string `json:"genres"`
	Languages   []string `json:"languages"`
	Title       string   `json:"title"`
	Size        string   `json:"size"`
	Difficulty  string   `json:"difficulty"`
	Status      []string `json:"status"`
	Modes       []string `json:"modes"`
	Description *string  `json:"description,omitempty"` // Champ facultatif
}

// edit playlist request
type PlaylistUpdateRequest struct {
	Title          string   `json:"title"`
	Genres         []string `json:"genres"`
	Languages      []string `json:"languages"`
	Tags           string   `json:"tags"`
	AddBeatmaps    []string `json:"addBeatmaps"`    // IDs des beatmaps à ajouter
	RemoveBeatmaps []string `json:"removeBeatmaps"` // IDs des beatmaps à retirer
}

/*
   -----------------
       HELPERS
   -----------------
*/

// helper pour calculer playlist genres, languages, modes, difficulty, length et size
func updatePlaylist(playlist *structures.Playlist) {
	modeSet := make(map[string]bool)
	languageSet := make(map[string]bool)
	genreSet := make(map[string]bool)

	playlist.Modes = []string{}
	playlist.Languages = []string{}
	playlist.Genres = []string{}

	for _, beatmap := range playlist.Beatmaps {
		modeSet[beatmap.Mode] = true
		languageSet[beatmap.Beatmapset.Language] = true
		genreSet[beatmap.Beatmapset.Genre] = true
	}

	for mode := range modeSet {
		if mode != "" {
			playlist.Modes = append(playlist.Modes, mode)
		}
	}
	for language := range languageSet {
		if language != "" {
			playlist.Languages = append(playlist.Languages, language)
		}
	}
	for genre := range genreSet {
		if genre != "" {
			playlist.Genres = append(playlist.Genres, genre)
		}
	}

	minDifficulty, maxDifficulty := diffRange(playlist.Beatmaps)
	playlist.Difficulty = fmt.Sprintf("%.2f - %.2f", minDifficulty, maxDifficulty)

	totalLength := 0
	for _, beatmap := range playlist.Beatmaps {
		totalLength += beatmap.Total_length
	}
	if len(playlist.Beatmaps) > 0 {
		playlist.Length = totalLength / len(playlist.Beatmaps)
	} else {
		playlist.Length = 0
	}

	playlist.Size = len(playlist.Beatmaps)
}

/*
   --------------------------------
       PLAYLISTS UPDATE HANDLERS
   --------------------------------
*/

// diff = min - max
func diffRange(beatmaps []structures.Beatmap) (minDifficulty, maxDifficulty float32) {
	if len(beatmaps) == 0 {
		return 0, 0
	}
	minDifficulty, maxDifficulty = beatmaps[0].Difficulty_rating, beatmaps[0].Difficulty_rating
	for _, beatmap := range beatmaps[1:] {
		if beatmap.Difficulty_rating < minDifficulty {
			minDifficulty = beatmap.Difficulty_rating
		}
		if beatmap.Difficulty_rating > maxDifficulty {
			maxDifficulty = beatmap.Difficulty_rating
		}
	}
	return minDifficulty, maxDifficulty
}

func nbBeatmaps(size string) int {
	switch size {
	case "petit":
		return 10
	case "moyen":
		return 30
	case "grand":
		return 60
	default:
		if customSize, err := strconv.Atoi(size); err == nil {
			return customSize
		}
		// Taille par défaut
		return 10
	}
}

// filter by genres, languages, tags, mode, numBeatmaps, length, difficulty, status
func filterBeatmaps(genres, languages, modes []string, numBeatmaps int, difficulty string, status []string) ([]structures.Beatmap, error) {
	var beatmaps []structures.Beatmap

	// filtre en fonction des préférences
	var filter []bson.M

	// si les slices sont vides, on ne les prend pas en compte
	if len(genres) > 0 && genres[0] != "" {
		filter = append(filter, bson.M{"beatmapset.genre": bson.M{"$in": genres}})
	}

	if len(languages) > 0 && languages[0] != "" {
		filter = append(filter, bson.M{"beatmapset.language": bson.M{"$in": languages}})
	}

	/*
	   Mode doit être l'un des suivants:
	       * osu
	       * taiko
	       * fruits
	       * mania
	*/
	fmt.Println("Modes : ", modes)
	if len(modes) > 0 && modes[0] != "" {

		filter = append(filter, bson.M{"mode": bson.M{"$in": modes}})
	}

	// plage de difficulté
	if difficulty != "" {
		diffParts := strings.Split(difficulty, "-")
		if len(diffParts) == 2 {
			minDiff, _ := strconv.ParseFloat(diffParts[0], 64)
			maxDiff, _ := strconv.ParseFloat(diffParts[1], 64)
			filter = append(filter, bson.M{"difficulty_rating": bson.M{"$gte": minDiff, "$lte": maxDiff}})
		}
	}

	/*
		Status doit être l'un des suivants:
			* ranked
			* approved
			* qualified
			* loved
			* pending
			* wip
			* graveyard
	*/
	if len(status) > 0 && status[0] != "" {
		filter = append(filter, bson.M{"status": bson.M{"$in": status}})
	}

	var filter2 bson.M
	if len(filter) > 0 {
		filter2 = bson.M{"$and": filter}
	} else {
		filter2 = bson.M{}
	}

	fmt.Println("Filtre : ", filter2)

	// trouver les beatmaps correspondantes
	cursor, err := data.BeatmapsCollection.Find(context.Background(), filter2)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// remplir le slice beatmaps avec les résultats
	for cursor.Next(context.Background()) {
		var beatmap structures.Beatmap
		if err := cursor.Decode(&beatmap); err != nil {
			return nil, err
		}
		beatmaps = append(beatmaps, beatmap)
	}

	if len(beatmaps) == 0 {
		return nil, errors.New("Aucune beatmap trouvée")
	}

	// Si le nombre de beatmaps trouvées est supérieur au nombre demandé, sélectionner aléatoirement
	if len(beatmaps) > numBeatmaps {
		selectedBeatmaps := make([]structures.Beatmap, 0, numBeatmaps)
		pickedIndexes := aux(len(beatmaps), numBeatmaps)
		for _, index := range pickedIndexes {
			selectedBeatmaps = append(selectedBeatmaps, beatmaps[index])
		}
		return selectedBeatmaps, nil
	}

	return beatmaps, nil
}

// aux aide à sélectionner des indices aléatoires d'un slice
func aux(sliceLen, numPicks int) []int {
	picked := make(map[int]bool)
	picks := make([]int, 0, numPicks)

	for len(picks) < numPicks {
		index := rand.Intn(sliceLen)
		if !picked[index] {
			picked[index] = true
			picks = append(picks, index)
		}
	}

	return picks
}

// create playlist
func HandleCreatePlaylist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Création de la playlist en cours...")

	//erreur méthode
	if r.Method != "POST" {
		http.Error(w, `{"status":"405","message":"Méthode non supportée","detail":"il faut une méthode POST"}`, http.StatusMethodNotAllowed)
		return
	}

	// vérifie si l'utilisateur est connecté
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, `{"status":"401","message":"Non autorisé: Aucun token de session fourni","detail":"l'utilisateur n'est pas connecté"}`, http.StatusUnauthorized)
		return
	}
	sessionToken := sessionCookie.Value

	sessionData, ok := Sessions[sessionToken]
	if !ok {
		http.Error(w, `{"status":"401","message":"Non autorisé: Session invalide ou expirée","detail":"la session a expiré"}`, http.StatusUnauthorized)
		return
	}
	fmt.Println("Session Data: ", sessionData.User_id)
	id, _ := primitive.ObjectIDFromHex(sessionData.User_id)

	var user structures.User
	if err := data.UsersCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user); err != nil {
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de la recherche de l'utilisateur `+sessionData.Username+` dans la base de données"}`, http.StatusInternalServerError)
		return
	}

	var request PlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"status":"400","message":"Données invalides","detail":"la requête n'est pas au format attendu"}`, http.StatusBadRequest)
		return
	}

	fmt.Println("Asked for ", request.Size, " beatmaps")

	//Difficulté sous forme: min-max
	numBeatmaps := nbBeatmaps(request.Size) //nombre de beatmaps voulu
	beatmaps, err := filterBeatmaps(request.Genres, request.Languages, request.Modes, numBeatmaps, request.Difficulty, request.Status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"nous n'avons pas trouvé de beatmaps"}`, http.StatusInternalServerError)
		return
	}

	fmt.Println("Found ", len(beatmaps), " beatmaps")
	fmt.Println("AVANT LE TRUC")
	// Traiter la requête
	description := ""
	if request.Description != nil {
		description = *request.Description // Déférencez le pointeur pour obtenir la valeur
	} else {
		generatedDescription, err := GeneratePlaylistDescription(request.Genres, request.Languages, request.Modes, request.Difficulty)
		if err != nil {
			description = "Description indisponible."
		} else {
			description = generatedDescription
		}
	}

	playlist_id := primitive.NewObjectID()
	var newPlaylist = structures.Playlist{
		Playlist_id: playlist_id,
		Title:       request.Title,
		Author:      user.Pseudo,
		Author_id:   sessionData.User_id,
		Url:         "http://osu.beatadvisor.com/playlist/" + playlist_id.Hex(),
		Beatmaps:    beatmaps,
		Size:        len(beatmaps),
		Genres:      request.Genres,
		Languages:   request.Languages,
		Modes:       request.Modes,
		Likes:       []string{},
		Comments:    []structures.Comment{},
		Description: description,
	}

	updatePlaylist(&newPlaylist)

	//erreur interne
	_, err = data.PlaylistsCollection.InsertOne(context.Background(), newPlaylist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de l'insertion de la playlist dans la base de données"}`, http.StatusInternalServerError)
		return
	}

	//on met à jour les playlists de l'utilisateur

	if err := data.UsersCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user); err != nil {
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de la recherche de l'utilisateur `+sessionData.Username+` dans la base de données"}`, http.StatusInternalServerError)
		return
	}

	update := bson.M{"$push": bson.M{"playlists": newPlaylist.Playlist_id}}
	if _, err := data.UsersCollection.UpdateOne(context.Background(), bson.M{"_id": id}, update); err != nil {
		http.Error(w, `{"status":"500","message":"Erreur interne lors de la mise à jour de l'utilisateur","detail":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	//création réussie
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "200",
		"message":  "Playlist créée avec succès",
		"playlist": newPlaylist,
	})
}

// get playlist
func HandleGetPlaylist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Récupération de la playlist en cours...")

	if r.Method != "GET" {
		http.Error(w, `{"status":"405","message":"Méthode non supportée","detail":"il faut une méthode GET"}`, http.StatusMethodNotAllowed)
		return
	}

	// vérifie si l'utilisateur est connecté
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, `{"status":"401","message":"Non autorisé: Aucun token de session fourni","detail":"l'utilisateur n'est pas connecté"}`, http.StatusUnauthorized)
		return
	}
	sessionToken := sessionCookie.Value

	_, ok := Sessions[sessionToken]
	if !ok {
		http.Error(w, `{"status":"401","message":"Non autorisé: Session invalide ou expirée","detail":"la session a expiré"}`, http.StatusUnauthorized)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		json.NewEncoder(w).Encode(Response{"400", "URL invalide", "url invalide"})
		http.Error(w, "URL invalide", http.StatusBadRequest)
		return
	}
	playlistID := pathParts[3]

	var playlist structures.Playlist
	objectId, err := primitive.ObjectIDFromHex(playlistID)
	if err != nil {
		http.Error(w, `{"status":"400","message":"ID de playlist invalide","detail":"l'ID de playlist fourni n'est pas valide"}`, http.StatusBadRequest)
		return
	}
	err = data.PlaylistsCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&playlist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, `{"status":"404","message":"Playlist introuvable","detail":"la playlist n'existe pas"}`, http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de la recherche de la playlist dans la base de données"}`, http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(playlist)
}

// delete playlist
func HandleDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Suppression de la playlist en cours...")

	if r.Method != "DELETE" {
		http.Error(w, `{"status":"405","message":"Méthode non supportée","detail":"il faut une méthode DELETE"}`, http.StatusMethodNotAllowed)
		return
	}

	// vérifie si l'utilisateur est connecté
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, `{"status":"401","message":"Non autorisé: Aucun token de session fourni","detail":"l'utilisateur n'est pas connecté"}`, http.StatusUnauthorized)
		return
	}
	sessionToken := sessionCookie.Value

	sessionData, ok := Sessions[sessionToken]
	if !ok {
		http.Error(w, `{"status":"401","message":"Non autorisé: Session invalide ou expirée","detail":"la session a expiré"}`, http.StatusUnauthorized)
		return
	}

	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 5 {
		http.Error(w, `{"status":"400","message":"URL invalide","detail":"url invalide"}`, http.StatusBadRequest)
		return
	}
	playlist_id := segments[4]

	id, err := primitive.ObjectIDFromHex(playlist_id)
	if err != nil {
		http.Error(w, `{"status":"400","message":"ID de playlist invalide","detail":"l'ID de playlist fourni n'est pas valide"}`, http.StatusBadRequest)
		return
	}

	fmt.Println("playlist_id: ", playlist_id)
	fmt.Println("ID: ", id)

	// vérifier si l'utilisateur est bien l'auteur
	var playlist structures.Playlist
	if err := data.PlaylistsCollection.FindOne(context.Background(), bson.M{"_id": id, "author_id": sessionData.User_id}).Decode(&playlist); err != nil {
		fmt.Println("Erreur: ", err)
		if err == mongo.ErrNoDocuments {
			http.Error(w, `{"status":"403","message":"Accès refusé","detail":"vous n'avez pas les droits pour modifier cette playlist ou elle n'existe pas"}`, http.StatusForbidden)
		} else {
			http.Error(w, `{"status":"500","message":"Erreur de serveur","detail":"erreur lors de la récupération de la playlist"}`, http.StatusInternalServerError)
		}
		return
	}

	result, err := data.PlaylistsCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de la suppression de la playlist `+playlist_id+` dans la base de données"}`, http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, `{"status":"404","message":"Playlist introuvable","detail":"la playlist `+playlist_id+` n'existe pas"}`, http.StatusNotFound)
		return
	}

	_, err = data.UsersCollection.UpdateMany(context.Background(), bson.M{}, bson.M{"$pull": bson.M{"Playlists": id}})
	if err != nil {
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de la mise à jour des utilisateurs"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Response{"200", "Playlist " + playlist_id + " supprimée avec succès", "playlist supprimée avec succès"})
}

// edit playlist
func HandleEditPlaylist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "PUT" {
		http.Error(w, `{"status":"405","message":"Méthode non supportée","detail":"il faut une méthode PUT"}`, http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 5 {
		http.Error(w, `{"status":"400","message":"URL invalide","detail":"url invalide"}`, http.StatusBadRequest)
		return
	}
	playlist_id := segments[4]

	fmt.Println("playlist_id: ", playlist_id)

	objectId, err := primitive.ObjectIDFromHex(playlist_id)
	if err != nil {
		http.Error(w, `{"status":"400","message":"ID de playlist invalide","detail":"l'ID de playlist fourni n'est pas valide"}`, http.StatusBadRequest)
		return
	}
	fmt.Println("Id = ", objectId)

	// Vérifie si l'utilisateur est connecté
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, `{"status":"401","message":"Non autorisé: Aucun token de session fourni","detail":"l'utilisateur n'est pas connecté"}`, http.StatusUnauthorized)
		return
	}
	sessionToken := sessionCookie.Value

	sessionData, ok := Sessions[sessionToken]
	if !ok {
		http.Error(w, `{"status":"401","message":"Non autorisé: Session invalide ou expirée","detail":"la session a expiré"}`, http.StatusUnauthorized)
		return
	}

	// vérifier si l'utilisateur est bien l'auteur
	var playlist structures.Playlist
	if err := data.PlaylistsCollection.FindOne(context.Background(), bson.M{"_id": objectId, "author_id": sessionData.User_id}).Decode(&playlist); err != nil {
		fmt.Println("Erreur: ", err)
		if err == mongo.ErrNoDocuments {
			http.Error(w, `{"status":"403","message":"Accès refusé","detail":"vous n'avez pas les droits pour modifier cette playlist ou elle n'existe pas"}`, http.StatusForbidden)
		} else {
			http.Error(w, `{"status":"500","message":"Erreur de serveur","detail":"erreur lors de la récupération de la playlist"}`, http.StatusInternalServerError)
		}
		return
	}

	// filtre
	var bupdates PlaylistUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&bupdates); err != nil {
		fmt.Println(bupdates)
		http.Error(w, `{"status":"400","message":"Données invalides","detail":"les données que vous avez envoyées sont invalides"}`, http.StatusBadRequest)
		return
	}

	update := bson.M{}
	if bupdates.Title != "" {
		update["$set"] = bson.M{"title": bupdates.Title}
	}
	if len(bupdates.Genres) > 0 {
		if update["$set"] == nil {
			update["$set"] = bson.M{}
		}
		update["$set"].(bson.M)["genres"] = bupdates.Genres
	}
	if len(bupdates.Languages) > 0 {
		if update["$set"] == nil {
			update["$set"] = bson.M{}
		}
		update["$set"].(bson.M)["languages"] = bupdates.Languages
	}
	if bupdates.Tags != "" {
		if update["$set"] == nil {
			update["$set"] = bson.M{}
		}
		update["$set"].(bson.M)["tags"] = []string{bupdates.Tags}
	}

	_, err = data.PlaylistsCollection.UpdateOne(context.Background(), bson.M{"_id": objectId}, update)
	if err != nil {
		http.Error(w, `{"status":"500","message":"Erreur lors de la mise à jour de la playlist","detail":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Response{"200", "Playlist mise à jour avec succès", "playlist mise à jour avec succès"})
}

// delete a beatmap from a playlist
func HandleDeleteBeatmapPlaylist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "DELETE" {
		http.Error(w, `{"status":"405","message":"Méthode non supportée","detail":"il faut une méthode DELETE"}`, http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, `{"status":"401","message":"Non autorisé: Aucun token de session fourni","detail":"l'utilisateur n'est pas connecté"}`, http.StatusUnauthorized)
		return
	}
	sessionToken := sessionCookie.Value

	sessionData, ok := Sessions[sessionToken]
	if !ok {
		http.Error(w, `{"status":"401","message":"Non autorisé: Session invalide ou expirée","detail":"la session a expiré"}`, http.StatusUnauthorized)
		return
	}

	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 7 {
		http.Error(w, `{"status":"400","message":"URL invalide","detail":"url invalide, ID de playlist et de beatmap attendus"}`, http.StatusBadRequest)
		return
	}
	playlist_id := segments[5]

	id, err := primitive.ObjectIDFromHex(playlist_id)
	if err != nil {
		http.Error(w, `{"status":"400","message":"ID de playlist invalide","detail":"l'id de playlist fourni n'est pas valide"}`, http.StatusBadRequest)
		return
	}

	beatmap_id := segments[6]

	beatmapID, err := strconv.Atoi(beatmap_id)
	if err != nil {
		http.Error(w, `{"status":"400","message":"ID de beatmap invalide","detail":"l'id beatmap fourni n'est pas valide"}`, http.StatusBadRequest)
		return
	}

	// vérifier si l'utilisateur est bien l'auteur
	var playlist structures.Playlist
	if err := data.PlaylistsCollection.FindOne(context.Background(), bson.M{"_id": id, "author_id": sessionData.User_id}).Decode(&playlist); err != nil {
		fmt.Println("Erreur: ", err)
		if err == mongo.ErrNoDocuments {
			http.Error(w, `{"status":"403","message":"Accès refusé","detail":"vous n'avez pas les droits pour modifier cette playlist ou elle n'existe pas"}`, http.StatusForbidden)
		} else {
			http.Error(w, `{"status":"500","message":"Erreur de serveur","detail":"erreur lors de la récupération de la playlist"}`, http.StatusInternalServerError)
		}
		return
	}

	update := bson.M{"$pull": bson.M{"beatmaps": bson.M{"id": beatmapID}}}
	_, err = data.PlaylistsCollection.UpdateOne(context.Background(), bson.M{"_id": id, "author_id": sessionData.User_id}, update)
	if err != nil {
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de la mise à jour de la playlist dans la base de données"}`, http.StatusInternalServerError)
		return
	}

	//update playlist après avoir supprimé la beatmap
	updatePlaylist(&playlist)

	update = bson.M{
		"$set": bson.M{
			"beatmaps":   playlist.Beatmaps,
			"modes":      playlist.Modes,
			"languages":  playlist.Languages,
			"genres":     playlist.Genres,
			"difficulty": playlist.Difficulty,
			"length":     playlist.Length,
			"size":       playlist.Size,
		},
	}

	_, err = data.PlaylistsCollection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, `{"status":"500","message":"Erreur interne lors de la mise à jour des détails de la playlist","detail":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Response{"200", "Beatmap supprimée avec succès de la playlist", "beatmap supprimée de la playlist avec succès"})
}

// add a playlist for a user
func HandleAddPlaylistForUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, `{"status": "405", "message": "Méthode non autorisée", "detail": "il faut une méthode POST"}`, http.StatusMethodNotAllowed)
		return
	}

	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, `{"status": "401", "message": "Non autorisé", "detail": "l'utilisateur n'est pas connecté"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"status": "500", "message": "Erreur interne", "detail": "erreur lors de la récupération du cookie de session"}`, http.StatusInternalServerError)
		return
	}

	sessionData, ok := Sessions[sessionToken.Value]
	if !ok {
		http.Error(w, `{ "status": "401", "message": "Non autorisé", "detail": "la session a expiré"}`, http.StatusUnauthorized)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, `{"status": "400", "message": "Erreur de traitement", "detail": "les données ne sont pas valides"}`, http.StatusBadRequest)
		return
	}

	playlistIDStr := r.FormValue("playlist_id")
	playlistID, err := primitive.ObjectIDFromHex(playlistIDStr)
	if err != nil {
		http.Error(w, `{"status": "400", "message": "ID de playlist invalide", "detail": "l'id de playlist fourni n'est pas valide"}`, http.StatusBadRequest)
		return
	}

	var playlist structures.Playlist
	if err := data.PlaylistsCollection.FindOne(context.Background(), bson.M{"_id": playlistID}).Decode(&playlist); err != nil {
		http.Error(w, `{"status": "404", "message": "Playlist introuvable", "detail": "la playlist `+playlistIDStr+` n'existe pas"}`, http.StatusNotFound)
		return
	}

	user_id, _ := primitive.ObjectIDFromHex(sessionData.User_id)
	filter := bson.M{"_id": user_id}
	update := bson.M{"$push": bson.M{"playlists": playlistID}}
	_, err = data.UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, `{"status": "500", "message": "Erreur interne", "detail": "erreur lors de l'ajout de la playlist à l'utilisateur"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Response{"200", "Playlist ajoutée avec succès à l'utilisateur", "playlist ajoutée avec succès à l'utilisateur " + sessionData.Username})
}

// get all playlists from a user
func HandleGetPlaylistsForUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, `{"status": "405", "message": "Méthode non autorisée", "detail": "il faut une méthode GET"}`, http.StatusMethodNotAllowed)
		return
	}

	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, `{"status": "401", "message": "Non autorisé", "detail": "l'utilisateur n'est pas connecté"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"status": "500", "message": "Erreur interne", "detail": "erreur lors de la récupération du cookie de session"}`, http.StatusInternalServerError)
		return
	}

	_, ok := Sessions[sessionToken.Value]
	if !ok {
		http.Error(w, `{ "status": "401", "message": "Non autorisé", "detail": "la session a expiré"}`, http.StatusUnauthorized)
		return
	}

	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 5 {
		http.Error(w, `{"status": "400", "message": "URL invalide", "detail": "url invalide"}`, http.StatusBadRequest)
		return
	}
	userIDStr := segments[4]

	fmt.Print("userIDStr: ", userIDStr)

	var user structures.User
	if err := data.UsersCollection.FindOne(context.Background(), bson.M{"pseudo": userIDStr}).Decode(&user); err != nil {
		http.Error(w, `{"status": "404", "message": "Utilisateur introuvable", "detail": "aucun utilisateur trouvé avec le pseudo spécifié"}`, http.StatusNotFound)
		return
	}

	playlists := make([]interface{}, 0, len(user.Playlists))
	for _, playlistID := range user.Playlists {
		var playlist structures.Playlist
		if err := data.PlaylistsCollection.FindOne(context.Background(), bson.M{"_id": playlistID}).Decode(&playlist); err == nil {
			playlists = append(playlists, playlist)
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "200",
		"message":   "Playlists récupérées avec succès",
		"playlists": playlists,
	})
}

func HandleGetAllPlaylists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, `{"status": "405", "message": "Méthode non autorisée", "detail": "il faut une méthode GET"}`, http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		json.NewEncoder(w).Encode(Response{"401", "Non autorisé: Aucun token de session fourni", "l'utilisateur n'est pas connecté"})
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}
	sessionToken := sessionCookie.Value

	_, ok := Sessions[sessionToken]
	if !ok {
		json.NewEncoder(w).Encode(Response{"401", "Non autorisé: Session invalide ou expirée", "la session a expiré"})
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	cursor, err := data.PlaylistsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, `{"status": "500", "message": "Erreur interne", "detail": "erreur lors de la récupération des playlists"}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	playlists := make([]map[string]interface{}, 0)
	for cursor.Next(context.Background()) {
		var playlist struct {
			ID    primitive.ObjectID `bson:"_id"`
			Title string             `bson:"title"`
		}
		if err := cursor.Decode(&playlist); err != nil {
			continue
		}
		// juste l'id et les titres
		playlists = append(playlists, map[string]interface{}{
			"Playlist_id": playlist.ID.Hex(),
			"Title":       playlist.Title,
		})
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, `{"status": "500", "message": "Erreur lors de l'itération sur les résultats", "detail": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "200",
		"message":   "Playlists récupérées avec succès",
		"playlists": playlists,
	})
}

const AIML_API_KEY = "49ba21f8ce234ea88db6c7c493dd38dc"
const AIML_API_URL = "https://api.aimlapi.com/v1/chat/completions"

// Fonction pour générer une description avec AIML API
func GeneratePlaylistDescription(genres, languages, modes []string, difficulty string) (string, error) {
	client := resty.New()
	fmt.Println("0")
	// Limiter les listes à 3 éléments max à cause des limites de l'api gratuite
	if len(genres) > 3 {
		genres = genres[:3]
	}
	if len(languages) > 3 {
		languages = languages[:3]
	}
	if len(modes) > 3 {
		modes = modes[:2]
	}
	// Construire le prompt
	prompt := fmt.Sprintf("Génère un bref paragraphe pour une playlist de beatmaps osu! en fonction de : Genres: %v, Langues: %v, Modes: %v, Difficulté: %s",
		genres, languages, modes, difficulty)

	// Construire la requête JSON
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]string{
			{"role": "system", "content": "Tu es un assistant expert d'osu!."},
			{"role": "user", "content": prompt},
		},
	})

	// Effectuer la requête POST
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+AIML_API_KEY).
		SetBody(requestBody).
		Post(AIML_API_URL)

	if err != nil {
		return "", fmt.Errorf("erreur lors de l'appel à l'API AIML: %v", err)
	}
	fmt.Println("1")
	fmt.Println("HTTP Status Code:", resp.StatusCode())
	fmt.Println("Raw Response:", resp.String())

	// Vérifier le statut HTTP
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("erreur HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	fmt.Println("2")

	// Décode la réponse JSON
	var response map[string]interface{}
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return "", fmt.Errorf("erreur lors du décodage JSON: %v", err)
	}
	fmt.Println("3")

	// Vérifier la présence de "choices"
	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("aucune réponse valide reçue de l'API AIML: %v", response)
	}
	fmt.Println("4")

	// Extraire la description
	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	if !ok {
		return "", fmt.Errorf("le format de réponse de l'API AIML ne correspond pas à ce qui était attendu")
	}

	fmt.Println("5")
	return message, nil
}
