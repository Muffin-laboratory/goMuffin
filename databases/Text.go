package databases

import (
	"time"

	"git.wh64.net/muffin/goMuffin/configs"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type InsertText struct {
	Text      string    `bson:"text" json:"text"`
	Persona   string    `bson:"persona" json:"persona"`
	CreatedAt time.Time `bson:"created_at"`
}

type Text struct {
	Id        bson.ObjectID `bson:"_id" json:"id"`
	Text      string        `bson:"text" json:"text"`
	Persona   string        `bson:"persona" json:"persona"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

var Texts *mongo.Collection = Client.Database(configs.Config.DBName).Collection("text")
