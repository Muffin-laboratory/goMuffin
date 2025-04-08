package databases

import (
	"time"

	"git.wh64.net/muffin/goMuffin/configs"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type InsertLearn struct {
	Command   string
	Result    string
	UserId    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
}

type Learn struct {
	Id        bson.ObjectID `bson:"_id"`
	Command   string
	Result    string
	UserId    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
}

var Learns *mongo.Collection = Client.Database(configs.Config.DBName).Collection("learn")
