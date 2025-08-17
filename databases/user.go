package databases

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ChattingMode int

type User struct {
	ID            bson.ObjectID `bson:"_id,omitempty"`
	UserID        string        `bson:"user_id,omitempty"`
	Blocked       bool          `bson:"blocked,omitempty"`
	BlockedReason string        `bson:"blocked_reason,omitempty"`
	ChatID        bson.ObjectID `bson:"chat_id,omitempty"`
	CreatedAt     time.Time     `bson:"created_at,omitempty"`
	ChattingMode  ChattingMode  `bson:"chatting_mode,omitempty"`
}

type UserCollection struct {
	*mongo.Collection
}

const (
	ChattingAIMode ChattingMode = iota + 1
	ChattingMuffinMode
)

func (c *UserCollection) IsUser(userId string) bool {
	var user *User
	c.FindOne(context.TODO(), User{UserID: userId}).Decode(&user)
	return user != nil
}

func (c *UserCollection) IsUserBlocked(userId string) (bool, string) {
	var user User
	err := c.FindOne(context.TODO(), User{UserID: userId}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, ""
		}

		fmt.Println(err)
		return true, "에러가 발생하여 차단한 유저를 구별 못해요. 계속 이러면 연락주세요."
	}
	return user.Blocked, user.BlockedReason
}

func (c *UserCollection) GetUserChattingMode(userId string) (ChattingMode, error) {
	var user User
	err := c.FindOne(context.TODO(), User{UserID: userId}).Decode(&user)
	if err != nil {
		return ChattingMuffinMode, err
	}
	return user.ChattingMode, nil
}

func ModeString(mode ChattingMode) string {
	switch mode {
	case ChattingAIMode:
		return "AI모드"
	case ChattingMuffinMode:
		return "머핀 모드"
	default:
		return "알 수 없음"
	}
}

func (u *User) ModeString() string {
	return ModeString(u.ChattingMode)
}
