package handlers

import (
	"beatadvisor/serveur/data"
	"beatadvisor/serveur/structures"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
   ----------------------------
       USERS UPDATE HANDLERS
   ----------------------------
*/

// login passé dans l'URL
func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		json.NewEncoder(w).Encode(Response{Status: "405", Message: "Méthode non supportée", Detail: "il faut une méthode GET"})
		http.Error(w, "Méthode non supportée", http.StatusMethodNotAllowed)
		return
	}

	// vérifie si l'utilisateur est connecté
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		json.NewEncoder(w).Encode(Response{Status: "401", Message: "Non autorisé: Aucun token de session fourni", Detail: "l'utilisateur n'est pas connecté"})
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}
	sessionToken := sessionCookie.Value

	_, ok := Sessions[sessionToken]
	if !ok {
		json.NewEncoder(w).Encode(Response{Status: "401", Message: "Non autorisé: Session invalide ou expirée", Detail: "la session a expiré"})
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 4 {
		json.NewEncoder(w).Encode(Response{Status: "400", Message: "URL invalide", Detail: "url invalide"})
		http.Error(w, "URL invalide", http.StatusBadRequest)
		return
	}
	pseudo := pathParts[3]
	fmt.Println(pseudo)

	var user structures.User
	err = data.UsersCollection.FindOne(context.Background(), bson.M{"pseudo": pseudo}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Response{Status: "404", Message: "Utilisateur non trouvé", Detail: "l'utilisateur n'existe pas"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Status: "500", Message: "Erreur interne", Detail: "erreur lors de la recherche de l'utilisateur dans la base de données"})
		}
		return
	}

	json.NewEncoder(w).Encode(user)
	fmt.Println("Utilisateur trouvé")
}

// delete User
func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "DELETE" {
		http.Error(w, `{"status":"405","message":"Méthode non supportée","detail":"Il faut une méthode GET"}`, http.StatusMethodNotAllowed)
		return
	}

	// vérifie si l'utilisateur est connecté
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		json.NewEncoder(w).Encode(Response{Status: "401", Message: "Non autorisé: Aucun token de session fourni", Detail: "l'utilisateur n'est pas connecté"})
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}
	sessionToken := sessionCookie.Value

	sessdata, ok := Sessions[sessionToken]
	if !ok {
		http.Error(w, `{"status":"401","message":"Non autorisé: Session invalide ou expirée","detail":"la session a expiré"}`, http.StatusUnauthorized)
		return
	}

	// ID, email ou pseudo
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 5 {
		http.Error(w, `{"status":"400","message":"URL invalide","detail":"url invalide"}`, http.StatusBadRequest)
		return
	}
	login := segments[4]

	id, err := primitive.ObjectIDFromHex(login) // ID
	var filter bson.M
	if err != nil { // Si ce n'est pas un ID valide, recherche par pseudo ou email
		filter = bson.M{"$or": []interface{}{
			bson.M{"email": login},
			bson.M{"pseudo": login},
		}}
	} else { // Si c'est un ID valide, recherche par _id
		filter = bson.M{"_id": id}
	}

	var user structures.User
	err = data.UsersCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, `{"status":"404","message":"Utilisateur introuvable","detail":"l'utilisateur `+sessdata.Username+` n'existe pas"}`, http.StatusNotFound)
		return
	}

	result, err := data.UsersCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de la suppression de l'utilisateur `+user.Pseudo+` dans la base de données"}`, http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, `{"status":"404","message":"Utilisateur introuvable","detail":"l'utilisateur `+user.Pseudo+` n'existe pas"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(Response{
		Status:  "200",
		Message: "Utilisateur supprimé avec succès. ID: " + user.User_id.Hex() + " Pseudo: " + user.Pseudo + ", Email: " + user.Email,
		Detail:  "utilisateur " + user.Pseudo + " supprimé avec succès",
	})
	fmt.Println("Utilisateur supprimé avec succès")
}
