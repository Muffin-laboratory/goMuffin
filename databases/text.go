package databases

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Text struct {
	Id        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Text      string        `bson:"text,omitempty" json:"text"`
	Persona   string        `bson:"persona,omitempty" json:"persona"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
}
