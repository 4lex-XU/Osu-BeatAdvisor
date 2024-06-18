package handlers

import (
	"beatadvisor/serveur/data"
	"beatadvisor/serveur/structures"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
   -----------------------
       SOCIAL FEATURES
   -----------------------
*/

// add comment to a playlist
func HandleAddComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Ajout de commentaire en cours...")

	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// vérifie si l'utilisateur est connecté
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, `{"status":"401","message":"Non autorisé: Aucun token de session fourni","detail":"l'utilisateur n'est pas connecté"}`, http.StatusUnauthorized)
		return
	}
	sessionToken := sessionCookie.Value

	fmt.Println("sessionToken add Comm: ", sessionToken)

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

	// Décode le commentaire reçu dans le corps de la requête
	var newComment structures.Comment
	if err := json.NewDecoder(r.Body).Decode(&newComment); err != nil {
		http.Error(w, `{"status":"400","message":"Données invalides","detail":"votre commentaire n'est pas valide"}`, http.StatusBadRequest)
		return
	}
	newComment.Id = primitive.NewObjectID()
	newComment.Author = sessionData.Username
	// Ajoute le commentaire à la playlist spécifiée
	updateResult, err := data.PlaylistsCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"comments": newComment}},
	)
	if err != nil {
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de l'ajout du commentaire dans la base de données"}`, http.StatusInternalServerError)
		return
	}

	if updateResult.ModifiedCount == 0 {
		http.Error(w, `{"status":"404","message":"Playlist introuvable","detail":"la playlist `+id.Hex()+` n'existe pas"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(Response{"200", "Commentaire ajouté avec succès", "commentaire ajouté avec succès, id = " + newComment.Id.Hex()})
}

// like a playlist
func HandleLike(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
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
		http.Error(w, `{"status":"400","message":"ID de playlist invalide","detail":"id de playlist invalide"}`, http.StatusBadRequest)
		return
	}

	// Incrémentation du compteur de likes pour la playlist spécifiée
	updateResult, err := data.PlaylistsCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$addToSet": bson.M{"likes": sessionData.Username}}, // un utilisateur ne peut liker qu'une fois
	)
	if err != nil {
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de l'ajout du like dans la base de données"}`, http.StatusInternalServerError)
		return
	}

	if updateResult.ModifiedCount == 0 {
		http.Error(w, `{"status":"404","message":"Playlist introuvable","detail":"la playlist `+id.Hex()+` n'existe pas"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(Response{"200", "Like ajouté avec succès", "vous avez likez la playlist"})
}

// delete like
func HandleDeleteLike(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, `{"status":"405","message":"Méthode non supportée","detail":"il faut une méthode POST"}`, http.StatusMethodNotAllowed)
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
	if len(segments) < 6 {
		http.Error(w, `{"status":"400","message":"URL invalide","detail":"url invalide"}`, http.StatusBadRequest)
		return
	}
	playlistID := segments[5]

	id, err := primitive.ObjectIDFromHex(playlistID)
	if err != nil {
		http.Error(w, `{"status":"400","message":"ID de playlist invalide","detail":"id de playlist invalide"}`, http.StatusBadRequest)
		return
	}

	updateResult, err := data.PlaylistsCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$pull": bson.M{"likes": sessionData.Username}},
	)
	if err != nil {
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de la suppression du like dans la base de données"}`, http.StatusInternalServerError)
		return
	}

	if updateResult.ModifiedCount == 0 {
		http.Error(w, `{"status":"404","message":"Playlist introuvable","detail":"la playlist `+playlistID+` n'existe pas"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(Response{"200", "Like supprimé avec succès", "like supprimé avec succès"})
}

// delete a comment from a playlist
func HandleDeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Suppression de commentaire en cours...")

	if r.Method != "DELETE" {
		http.Error(w, `{"status":"405","message":"Méthode non supportée","detail":"il faut une méthode DELETE"}`, http.StatusMethodNotAllowed)
		return
	}

	// Vérifie si l'utilisateur est connecté
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

	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 7 {
		http.Error(w, `{"status":"400","message":"URL invalide","detail":"url invalide"}`, http.StatusBadRequest)
		return
	}
	playlistIDStr := segments[5]
	commentIDStr := segments[6]

	fmt.Println("playlistIDStr: ", playlistIDStr)
	fmt.Println("commentIDStr: ", commentIDStr)

	playlistID, err := primitive.ObjectIDFromHex(playlistIDStr)
	if err != nil {
		http.Error(w, `{"status":"400","message":"ID de playlist invalide","detail":"id de playlist invalide"}`, http.StatusBadRequest)
		return
	}

	commentID, err := primitive.ObjectIDFromHex(commentIDStr)
	if err != nil {
		http.Error(w, `{"status":"400","message":"ID de commentaire invalide","detail":"id du commentaire invalide"}`, http.StatusBadRequest)
		return
	}

	updateResult, err := data.PlaylistsCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": playlistID},
		bson.M{"$pull": bson.M{"comments": bson.M{"_id": commentID}}},
	)
	if err != nil {
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de la suppression du commentaire `+commentIDStr+` dans la base de données"}`, http.StatusInternalServerError)
		return
	}

	if updateResult.ModifiedCount == 0 {
		http.Error(w, `{"status":"404","message":"Commentaire introuvable","detail":"le commentaire `+commentIDStr+` ou la playlist `+playlistIDStr+` n'existe pas"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(Response{"200", "Commentaire supprimé avec succès", "commentaire supprimé avec succès"})
}

// share a playlist
func HandleSharePlaylist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "Méthode non supportée", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 5 {
		json.NewEncoder(w).Encode(Response{"400", "URL invalide", "url invalide"})
		http.Error(w, "URL invalide", http.StatusBadRequest)
		return
	}
	playlist_id := segments[4]

	id, err := primitive.ObjectIDFromHex(playlist_id)
	if err != nil {
		http.Error(w, `{"status":"400","message":"ID de playlist invalide","detail":"l'ID de playlist fourni n'est pas valide"}`, http.StatusBadRequest)
		return
	}

	var playlist structures.Playlist
	err = data.PlaylistsCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&playlist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, `{"status":"404","message":"Playlist introuvable","detail":"la playlist n'existe pas"}`, http.StatusNotFound)
			return
		} else {
			http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de la recherche de la playlist dans la base de données"}`, http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "200",
		"message": "Playlist partagée avec succès",
		"link":    playlist.Url,
	})
}

// profile editing
func HandleEditProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "PUT" {
		http.Error(w, `{"status":"405","message":"Méthode non supportée","detail":"il faut une méthode PUT"}`, http.StatusMethodNotAllowed)
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
	fmt.Println("sessionData: ", sessionData)

	if err := r.ParseForm(); err != nil {
		http.Error(w, `{"status":"400","message":"Données invalides","detail":"les données que vous avez envoyées sont invalides"}`, http.StatusBadRequest)
		return
	}

	login := r.FormValue("login")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")
	naissance := r.FormValue("naissance")
	description := r.FormValue("description")
	ville := r.FormValue("ville")

	update := bson.M{}
	if login != "" && login != sessionData.Username {
		// Vérifie si le pseudo ou l'email est déjà utilisé

		var user structures.User
		err := data.UsersCollection.FindOne(context.Background(), bson.M{"pseudo": login}).Decode(&user)
		if err != nil && err != mongo.ErrNoDocuments {
			http.Error(w, `{"status":500,"message":"Erreur interne","detail":"erreur lors de la vérification de l'unicité du login dans la base de données"}`, http.StatusInternalServerError)
			return
		}

		update["pseudo"] = login
		sessionData := Sessions[sessionToken]
		sessionData.Username = login // Mise à jour de la session
	}

	if password != "" && confirmPassword != "" && password == confirmPassword {
		update["password"] = password
	} else if password != confirmPassword {
		http.Error(w, `{"status":"400","message":"Mots de passe non identiques","detail":"les mots de passe ne correspondent pas"}`, http.StatusBadRequest)
		return
	}

	if naissance != "" {
		if _, err := time.Parse("2006-01-02", naissance); err == nil {
			update["naissance"] = naissance
		} else {
			http.Error(w, `{"status":"400","message":"Date invalide","detail":"la date de naissance doit être au format 'YYYY-MM-DD'"}`, http.StatusBadRequest)
			return
		}
	}

	if description != "" {
		update["description"] = description
	}

	if ville != "" {
		update["ville"] = ville
	}

	if len(update) > 0 {
		id, _ := primitive.ObjectIDFromHex(sessionData.User_id)
		if _, err := data.UsersCollection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": update}); err != nil {
			http.Error(w, `{"status":"500","message":"Erreur interne lors de la mise à jour du profil","detail":"erreur lors de la mise à jour du profil dans la base de données"}`, http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(Response{"200", "Profil mis à jour avec succès", "profil mis à jour avec succès"})
}

// add a friend (follow)
func HandleAddFriend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "PUT" {
		http.Error(w, `{"status":"405","message":"Méthode non supportée","detail":"il faut une méthode PUT"}`, http.StatusMethodNotAllowed)
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
	id, _ := primitive.ObjectIDFromHex(sessionData.User_id)

	r.ParseForm()
	friend := r.FormValue("friend")

	filter := bson.M{"_id": id}
	update := bson.M{"$addToSet": bson.M{"friends": friend}}

	result, err := data.UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de la mise à jour de la liste d'amis"}`, http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, `{"status":"404","message":"Aucune mise à jour effectuée","detail":"l'ami est peut-être déjà dans la liste ou l'utilisateur n'existe pas"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(Response{"200", "Ami ajouté avec succès", "vous suivez cet utilisateur désormais"})
}

// check if a user is following another user
func HandleCheckFollow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, `{"status":"405","message":"Méthode non supportée","detail":"il faut une méthode GET"}`, http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 5 {
		http.Error(w, `{"status":"400","message":"URL invalide","detail":"url invalide"}`, http.StatusBadRequest)
		return
	}
	friend := segments[4]

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
	id, _ := primitive.ObjectIDFromHex(sessionData.User_id)

	var user structures.User
	err = data.UsersCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de la récupération des données de l'utilisateur"}`, http.StatusInternalServerError)
		return
	}

	isFriend := false
	for _, friendPseudo := range user.Friends {
		if friendPseudo == friend {
			isFriend = true
			break
		}
	}

	json.NewEncoder(w).Encode(map[string]bool{"isFriend": isFriend})
}

// unfollow
func HandleRemoveFriend(w http.ResponseWriter, r *http.Request) {
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
	id, _ := primitive.ObjectIDFromHex(sessionData.User_id)

	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 6 {
		http.Error(w, `{"status":"400","message":"URL invalide", "detail":"url invalide"}`, http.StatusBadRequest)
		return
	}
	friendID := segments[5]

	filter := bson.M{"_id": id}
	update := bson.M{"$pull": bson.M{"friends": friendID}}

	result, err := data.UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, `{"status":"500","message":"Erreur interne","detail":"erreur lors de la mise à jour de la liste d'amis"}`, http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, `{"status":"404","message":"Aucune mise à jour effectuée","detail":"l'ami n'est peut-être pas dans la liste ou l'utilisateur n'existe pas"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(Response{"200", "Ami retiré avec succès", "vous ne suivez plus cet utilisateur"})
}
