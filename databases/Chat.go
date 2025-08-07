package databases

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Chat struct {
	Id        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name,omitempty"`
	UserId    string        `bson:"user_id,omitempty"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
}

func CreateChat(userId, name string) (*mongo.InsertOneResult, error) {
	createdChat, err := Database.Chats.InsertOne(context.TODO(), Chat{UserId: userId, Name: name, CreatedAt: time.Now()})
	if err != nil {
		return nil, err
	}

	_, err = Database.Users.UpdateOne(context.TODO(), User{UserId: userId}, bson.D{{
		Key:   "$set",
		Value: User{ChatId: createdChat.InsertedID.(bson.ObjectID)},
	}})
	if err != nil {
		return nil, err
	}
	return createdChat, nil
}
