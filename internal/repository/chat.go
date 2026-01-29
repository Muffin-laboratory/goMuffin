package repository

import (
	"context"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Chat struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name,omitempty"`
	UserID    string        `bson:"user_id,omitempty"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
}

type ChatCollection struct {
	coll *mongo.Collection
}

func (c *ChatCollection) Create(ctx context.Context, userID, name string) (*mongo.InsertOneResult, error) {
	createdChat, err := c.coll.InsertOne(ctx, Chat{UserID: userID, Name: name, CreatedAt: time.Now()})
	if err != nil {
		return nil, err
	}

	chatID := createdChat.InsertedID.(bson.ObjectID)
	if _, err := GetDatabase().Users.Update(ctx, userID, &UserUpdate{
		ChatID: &chatID,
	}); err != nil {
		return nil, err
	}

	return createdChat, nil
}

func (c *ChatCollection) Find(ctx context.Context, filter query.QueryBuilder) ([]Chat, error) {
	cur, err := c.coll.Find(ctx, filter.Build())
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var chats []Chat

	if err := cur.All(ctx, &chats); err != nil {
		return nil, err
	}

	return chats, nil
}

func (c *ChatCollection) FindOne(ctx context.Context, filter query.QueryBuilder) (*Chat, error) {
	var chat Chat

	if err := c.coll.FindOne(ctx, filter.Build()).Decode(&chat); err != nil {
		return nil, err
	}

	return &chat, nil
}

func (c *ChatCollection) FindByID(ctx context.Context, id bson.ObjectID) (*Chat, error) {
	return c.FindOne(ctx, query.ChatQueryBuilder().SetID(id))
}

func (c *ChatCollection) DeleteByID(ctx context.Context, id bson.ObjectID) error {
	_, err := c.coll.DeleteOne(ctx, query.ChatQueryBuilder().SetID(id))
	return err
}
