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

func (c *ChatCollection) CreateChat(userID, name string) (*mongo.InsertOneResult, error) {
	createdChat, err := c.InsertOne(context.TODO(), Chat{UserID: userID, Name: name, CreatedAt: time.Now()})
	if err != nil {
		return nil, err
	}

	if _, err := GetDatabase().Users.Update(userID, User{
		ChatID: createdChat.InsertedID.(bson.ObjectID),
	}); err != nil {
		return nil, err
	}

	return createdChat, nil
}
