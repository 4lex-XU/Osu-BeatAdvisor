package structures

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	Id     primitive.ObjectID `bson:"_id"`
	Author string             `bson:"author" json:"author"`
	Text   string             `bson:"text" json:"text"`
	Date   string         		`bson:"date" json:"date"`
}
