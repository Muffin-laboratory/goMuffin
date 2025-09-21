package databases

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Knowledge struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Command   string        `bson:"command,omitempty"`
	Result    string        `bson:"result,omitempty"`
	UserID    string        `bson:"user_id,omitempty"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
}
