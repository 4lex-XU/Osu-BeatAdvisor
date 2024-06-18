package handlers

import (
	"beatadvisor/serveur/data"
	"beatadvisor/serveur/structures"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
   ------------------------------------
       Response/Requests Structures
   ------------------------------------
*/

type AuthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
	UserId  string `json:"user_id"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

var (
	dbName = data.DbName

	Sessions = make(map[string]structures.SessionData)
)

/*
   --------------------------------
       AUTHENTIFICATION HANDLERS
   --------------------------------
*/

// register
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Inscription en cours")

	if r.Method != "PUT" {
		http.Error(w, `{"status":405,"message":"Méthode non supportée","detail":"il faut une méthode PUT"}`, http.StatusMethodNotAllowed)

		return
	}

	r.ParseForm()
	email := r.FormValue("email")
	pseudo := r.FormValue("pseudo")
	password := r.FormValue("password")

	fmt.Println("Email: ", email)
	fmt.Println("Pseudo: ", pseudo)
	fmt.Println("Password: ", password)

	//email ou pseudo ou password vides
	if email == "" || password == "" || pseudo == "" {

		http.Error(w, `{"status":400,"message":"Champs manquants","detail":"les champs email, pseudo et password sont nécessaires"}`, http.StatusBadRequest)
		return
	}

	//vérifie si l'email est valide (@ et .com/.fr/etc)
	regexEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !regexEmail.MatchString(email) {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, `{"status":400,"message":"Email invalide","detail":"l'email doit contenir @ et finir par .com/.fr/..."}`, http.StatusBadRequest)

		return
	}

	// l'utilisateur existe déjà (pseudo et email)
	var existingUser structures.User
	_ = data.UsersCollection.FindOne(context.Background(), bson.M{
		"$or": []bson.M{
			{"email": email},
			{"pseudo": pseudo},
		},
	}).Decode(&existingUser)

	//fmt.Println("Existing User: ", existingUser)
	if existingUser.Pseudo != "" || existingUser.Email != "" {
		http.Error(w, `{"status":409,"message":"Utilisateur existant","detail":"l'utilisateur `+existingUser.Pseudo+` existe déjà ou email déjà utilisé"}`, http.StatusConflict)
		return
	}

	newUser := structures.NewUser(email, pseudo, password)

	_, err := data.UsersCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		http.Error(w, `{"status":500,"message":"Erreur interne","detail":"erreur lors de l'insertion de l'utilisateur `+newUser.Pseudo+` dans la base de données"}`, http.StatusInternalServerError)
		return
	}

	// Inscription réussie
	json.NewEncoder(w).Encode(Response{Status: "200", Message: "Utilisateur inscrit avec succès", Detail: "inscription réussie pour " + pseudo})
}

// login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, `{"status":"405","message":"Méthode non supportée","detail":"il faut une méthode POST"}`, http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	// login = email ou pseudo
	login := r.FormValue("login")
	password := r.FormValue("password")

	//login ou password vides
	if login == "" || password == "" {
		http.Error(w, `{"status":400,"message":"Champs manquants","detail":"les champs login et password sont nécessaires"}`, http.StatusBadRequest)
		return
	}

	var user structures.User
	err := data.UsersCollection.FindOne(context.Background(), bson.M{"$or": []interface{}{
		bson.M{"email": login},
		bson.M{"pseudo": login},
	}, "password": password}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, `{"status":"403","message":"Login et/ou mot de passe invalide(s)","detail":"utilisateur inconnu ou mot de passe incorrect"}`, http.StatusUnauthorized)
		} else {
			fmt.Println("Erreur lors de la recherche de l'utilisateur:", err)
			http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de la recherche de l'utilisateur `+login+` dans la base de données"}`, http.StatusInternalServerError)
		}
		return
	}

	// Session
	sessionToken := user.User_id.Hex()
	exptime := time.Now().Add(1 * time.Hour) // modifier si besoin

	Sessions[sessionToken] = structures.SessionData{
		User_id:  user.User_id.Hex(),
		Username: user.Pseudo,
		Expire:   exptime,
	}

	fmt.Println("Session Token: ", sessionToken)

	// Créez un cookie pour la session
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  exptime,
	})

	// réussie
	response := AuthResponse{
		Status:  "200",
		Message: "Connexion réussie",
		Detail:  "bienvenue " + user.Pseudo,
		UserId:  user.User_id.Hex(),
	}
	json.NewEncoder(w).Encode(response)
	fmt.Fprintf(w, "Connexion réussie pour %s", user.Pseudo)
}

// log out
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, `{"status":"401","message":"Pas de cookie de session trouvé","detail":"l'utilisateur n'est pas connecté"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"status":"500","message":"Erreur interne du serveur","detail":"impossible de récupérer le cookie de session"}`, http.StatusInternalServerError)
		return
	}

	// Supprimez la session du stockage de session
	delete(Sessions, sessionToken.Value)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	response := AuthResponse{
		Status:  "200",
		Message: "Déconnexion réussie",
		UserId:  sessionToken.Value,
	}

	json.NewEncoder(w).Encode(response)
}
