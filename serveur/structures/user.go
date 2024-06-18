package structures

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	User_id     primitive.ObjectID `bson:"_id"`
	Email       string
	Pseudo      string
	Password    string
	Playlists   []primitive.ObjectID `bson:"playlists"`
	Friends     []string
	Naissance   string
	Description string
	Ville       string
}

type SessionData struct {
	User_id  string
	Username string
	Expire   time.Time
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

func NewUser(email, pseudo, password string) *User {
	// lock.Lock()
	// defer lock.Unlock()
	// cpt++
	return &User{
		User_id:   primitive.NewObjectID(),
		Email:     email,
		Pseudo:    pseudo,
		Password:  password,
		Playlists: []primitive.ObjectID{},
		Friends:   []string{},
	}
}

/*
   -----------------
       Methods
   -----------------
*/

// func checkPassword(client user, login string, password string) bool {
// 	//TODO: utile?
// }
