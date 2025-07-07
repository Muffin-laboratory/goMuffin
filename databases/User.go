package databases

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	Id            bson.ObjectID `bson:"_id,omitempty"`
	UserId        string        `bson:"user_id"`
	Blocked       bool          `bson:"blocked"`
	BlockedReason string        `bson:"blocked_reason"`
	CreatedAt     time.Time     `bson:"created_at"`
}

func (d *MuffinDatabase) IsUser(userId string) bool {
	var user *User
	d.Users.FindOne(context.TODO(), bson.D{{Key: "user_id", Value: userId}}).Decode(&user)
	return user != nil
}

func (d *MuffinDatabase) IsUserBlocked(userId string) (bool, string) {
	var user User
	d.Users.FindOne(context.TODO(), bson.D{{Key: "user_id", Value: userId}}).Decode(&user)
	return user.Blocked, user.BlockedReason
}
