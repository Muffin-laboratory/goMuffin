package databases

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Learn struct {
	Id        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Command   string        `bson:"command"`
	Result    string        `bson:"result"`
	UserId    string        `bson:"user_id"`
	CreatedAt time.Time     `bson:"created_at"`
}
