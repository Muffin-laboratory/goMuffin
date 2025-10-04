package databases

import (
	"context"
	"time"

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
	*mongo.Collection
}

func (c *ChatCollection) Create(userID, name string) (*mongo.InsertOneResult, error) {
	createdChat, err := c.InsertOne(context.TODO(), Chat{UserID: userID, Name: name, CreatedAt: time.Now()})
	if err != nil {
		return nil, err
	}

	chatID := createdChat.InsertedID.(bson.ObjectID)
	if _, err := GetDatabase().Users.Update(userID, &UserUpdate{
		ChatID: &chatID,
	}); err != nil {
		return nil, err
	}

	return createdChat, nil
}
