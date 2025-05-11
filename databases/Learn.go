package databases

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type InsertLearn struct {
	Command   string    `bson:"command"`
	Result    string    `bson:"result"`
	UserId    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
}

type Learn struct {
	Id        bson.ObjectID `bson:"_id" json:"id"`
	Command   string        `bson:"command" json:"command"`
	Result    string        `bson:"result" json:"result"`
	UserId    string        `bson:"user_id" json:"user_id"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}
