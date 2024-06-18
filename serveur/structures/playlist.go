package structures

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Playlist struct {
	Playlist_id primitive.ObjectID `bson:"_id"`
	Title       string
	Author      string
	Author_id   string
	Url         string
	Length      int
	Modes       []string
	Genres      []string
	Languages   []string
	Tags        []string
	Difficulty  string

	Beatmaps []Beatmap
	Size     int

	Likes    []string
	Comments []Comment
}

type Beatmap struct {
	Difficulty_rating float32
	Id                int
	Mode              string
	Status            string
	Total_length      int
	User_id           int
	Url               string
	Beatmapset        Beatmapset
}

type Beatmapset struct {
	Title          string
	Genre          string
	Language       string
	Tags           string
	Ranked_date    string
	Submitted_date string
}

// var (
//     cpt int = 0
//     lock          sync.Mutex
// )

/*
   -------------------
       Constructor
   -------------------
*/

func NewPlaylist(title string, url string, modes []string, difficulty string, genres, languages, tags []string, beatmaps []Beatmap, size int) *Playlist {
	// lock.Lock()
	// defer lock.Unlock()
	// userIDCounter++
	return &Playlist{
		//playlist_id: cpt,
		Playlist_id: primitive.NewObjectID(),
		Title:       title,
		Url:         url,
		Modes:       modes,
		Difficulty:  difficulty,
		Genres:      genres,
		Languages:   languages,
		Tags:        tags,
		Beatmaps:    beatmaps,
		Size:        size,
	}
}
