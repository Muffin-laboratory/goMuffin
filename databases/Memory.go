package databases

import "go.mongodb.org/mongo-driver/v2/bson"

type Memory struct {
	Id      bson.ObjectID `bson:"_id,omitempty"`
	UserId  string        `bson:"user_id,omitempty"`
	Content string        `bson:"content,omitempty"`
	Answer  string        `bson:"answer,omitempty"`
	ChatId  bson.ObjectID `bson:"chat_id,omitempty"`
}
