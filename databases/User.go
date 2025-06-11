package databases

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type InsertUser struct {
	UserId    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
}

type User struct {
	Id        bson.ObjectID `bson:"_id"`
	UserId    string        `bson:"user_id"`
	CreatedAt time.Time     `bson:"created_at"`
}

func (d *MuffinDatabase) IsUser(userId string) bool {
	var user *User
	d.Users.FindOne(context.TODO(), bson.D{{Key: "user_id", Value: userId}}).Decode(&user)

	if user != nil {
		return true
	}
	return false
}
