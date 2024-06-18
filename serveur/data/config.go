package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// Base de données
	Uri = "mongodb+srv://luluantoinex:324CfsMnNT9L7tfg@beatadvisor.uw3ukhu.mongodb.net/"

	// Client
	Client *mongo.Client

	// Nom de la bd
	DbName = "BeatAdvisor"

	// Users collection dans la bd
	UsersCollection *mongo.Collection
	// Playlists
	PlaylistsCollection *mongo.Collection
	// Beatmaps
	BeatmapsCollection *mongo.Collection
)

func Init() {
	var err error
	// Initialisation du client MongoDB
	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(Uri))
	if err != nil {
		log.Fatal(err)
	}

	// Vérification de la connexion
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := Client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	UsersCollection = Client.Database(DbName).Collection("Users")
	PlaylistsCollection = Client.Database(DbName).Collection("Playlists")
	BeatmapsCollection = Client.Database(DbName).Collection("Beatmaps")
}
