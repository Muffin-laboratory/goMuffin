package databases

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Memory struct {
	ID      bson.ObjectID `bson:"_id,omitempty"`
	UserID  string        `bson:"user_id,omitempty"`
	Content string        `bson:"content,omitempty"`
	Answer  string        `bson:"answer,omitempty"`
	ChatID  bson.ObjectID `bson:"chat_id,omitempty"`
}

type MemoryCollection struct {
	Collection *mongo.Collection
}

func (c *MemoryCollection) Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
	return c.Collection.Find(ctx, filter, opts...)
}
