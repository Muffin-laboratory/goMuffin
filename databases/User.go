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
	Id            bson.ObjectID `bson:"_id,omitempty"`
	UserId        string        `bson:"user_id,omitempty"`
	Blocked       bool          `bson:"blocked,omitempty"`
	BlockedReason string        `bson:"blocked_reason,omitempty"`
	ChatId        bson.ObjectID `bson:"chat_id,omitempty"`
	CreatedAt     time.Time     `bson:"created_at,omitempty"`
	ChattingMode  ChattingMode  `bson:"chatting_mode,omitempty"`
}

const (
	ChattingAIMode ChattingMode = iota + 1
	ChattingMuffinMode
)

func (d *MuffinDatabase) IsUser(userId string) bool {
	var user *User
	d.Users.FindOne(context.TODO(), User{UserId: userId}).Decode(&user)
	return user != nil
}

func (d *MuffinDatabase) IsUserBlocked(userId string) (bool, string) {
	var user User
	err := d.Users.FindOne(context.TODO(), User{UserId: userId}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, ""
		}

		fmt.Println(err)
		return true, "에러가 발생하여 차단한 유저를 구별 못해요. 계속 이러면 연락주세요."
	}
	return user.Blocked, user.BlockedReason
}

func (d *MuffinDatabase) GetUserChattingMode(userId string) (ChattingMode, error) {
	var user User
	err := d.Users.FindOne(context.TODO(), User{UserId: userId}).Decode(&user)
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
