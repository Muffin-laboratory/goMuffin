package databases

import "go.mongodb.org/mongo-driver/v2/bson"

type Memory struct {
	Id      bson.ObjectID `bson:"_id,omitempty"`
	UserId  string        `bson:"user_id"`
	Content string        `bson:"content"`
	Answer  string        `bson:"answer"`
}
