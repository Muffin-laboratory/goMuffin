package repository

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Text struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Text      string        `bson:"text,omitempty"`
	Persona   string        `bson:"persona,omitempty"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
}

type TextCollection struct {
	coll   *mongo.Collection
	caches struct {
		data []Text
		mu   *sync.RWMutex
	}
}

func newTextCollection(coll *mongo.Collection) *TextCollection {
	return &TextCollection{
		coll: coll,
		caches: struct {
			data []Text
			mu   *sync.RWMutex
		}{
			mu: &sync.RWMutex{},
		},
	}
}

func (c *TextCollection) Create(ctx context.Context, text string) (*Text, error) {
	data := Text{
		Text:      text,
		Persona:   "muffin",
		CreatedAt: time.Now(),
	}

	result, err := c.coll.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	data.ID = result.InsertedID.(bson.ObjectID)

	c.caches.mu.Lock()
	c.caches.data = append(c.caches.data, data)
	c.caches.mu.Unlock()

	return &data, nil
}

func (c *TextCollection) All(ctx context.Context) ([]Text, error) {
	if len(c.caches.data) > 0 {
		c.caches.mu.RLock()
		data := make([]Text, len(c.caches.data))
		copy(data, c.caches.data)
		c.caches.mu.RUnlock()
		return data, nil
	}

	cur, err := c.coll.Find(ctx, bson.D{{Key: "persona", Value: "muffin"}})
	if err != nil {
		return nil, err
	}

	var data []Text

	defer func() {
		ctx := context.WithoutCancel(ctx)
		if err := cur.Close(ctx); err != nil {
			slog.Error("failed to close text cursor", "error", err)
		}
	}()

	if err = cur.All(ctx, &data); err != nil {
		return nil, err
	}

	c.caches.mu.Lock()
	c.caches.data = append(c.caches.data, data...)
	c.caches.mu.Unlock()

	return data, nil
}
