package databases

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Text struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Text      string        `bson:"text,omitempty"`
	Persona   string        `bson:"persona,omitempty"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
}
