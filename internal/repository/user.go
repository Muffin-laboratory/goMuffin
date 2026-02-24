package repository

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/cache"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ChattingMode int

type User struct {
	ID                        int64         `bson:"_id,omitempty"`
	Blocked                   bool          `bson:"blocked"`
	BlockedReason             string        `bson:"blocked_reason"`
	ChatID                    bson.ObjectID `bson:"chat_id"`
	CreatedAt                 time.Time     `bson:"created_at"`
	ChattingMode              ChattingMode  `bson:"chatting_mode"`
	ReplyUser                 bool          `bson:"reply_user"`
	CreateNewChatAfter12Hours bool          `bson:"create_new_chat_after_12_hours"`
	Prompt                    string        `bson:"prompt"`
}

type UserUpdate struct {
	Blocked                   *bool          `bson:"blocked,omitempty"`
	BlockedReason             *string        `bson:"blocked_reason,omitempty"`
	ChatID                    *bson.ObjectID `bson:"chat_id,omitempty"`
	CreatedAt                 *time.Time     `bson:"created_at,omitempty"`
	ChattingMode              *ChattingMode  `bson:"chatting_mode,omitempty"`
	ReplyUser                 *bool          `bson:"reply_user,omitempty"`
	CreateNewChatAfter12Hours *bool          `bson:"create_new_chat_after_12_hours,omitempty"`
	Prompt                    *string        `bson:"prompt,omitempty"`
}

const (
	ChattingAIMode ChattingMode = iota + 1
	ChattingMuffinMode
)

type UserCollection struct {
	coll   *mongo.Collection
	caches *cache.CacheManager[int64, User]
}

func newUserCollection(coll *mongo.Collection) *UserCollection {
	return &UserCollection{
		coll:   coll,
		caches: cache.New[int64, User](timeToExpire),
	}
}

func (c *UserCollection) Create(ctx context.Context, userID int64) (*User, error) {
	user := User{
		ID:           userID,
		ChattingMode: ChattingAIMode,
		ReplyUser:    true,
		CreatedAt:    time.Now(),
	}

	_, err := c.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	c.caches.Set(userID, user)

	return &user, nil
}

func (c *UserCollection) All(ctx context.Context) ([]User, error) {
	if caches := c.caches.All(); len(caches) != 0 {
		data := make([]User, len(caches))

		for _, user := range caches {
			data = append(data, user)
		}

		return data, nil
	}

	var data []User

	cur, err := c.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer func() {
		ctx := context.WithoutCancel(ctx)
		if err := cur.Close(ctx); err != nil {
			slog.Error("failed to close user cursor", "error", err)
		}
	}()

	if err = cur.All(ctx, &data); err != nil {
		return nil, err
	}

	for _, data := range data {
		c.caches.Set(data.ID, data)
	}

	return data, nil
}

func (c *UserCollection) FindBlockedUser(ctx context.Context) ([]User, error) {
	var data []User

	cur, err := c.coll.Find(ctx, bson.D{{Key: "blocked", Value: true}})
	if err != nil {
		return nil, err
	}

	defer func() {
		ctx := context.WithoutCancel(ctx)
		if err := cur.Close(ctx); err != nil {
			slog.Error("failed to close user cursor", "error", err)
		}
	}()

	if err = cur.All(ctx, &data); err != nil {
		return nil, err
	}

	return data, err
}

func (c *UserCollection) FindByID(ctx context.Context, userID int64) (*User, error) {
	if user, ok := c.caches.Get(userID); ok {
		return &user, nil
	}

	var user User

	if err := c.coll.FindOne(ctx, bson.M{
		"_id": userID,
	}).Decode(&user); err != nil {
		return nil, err
	}

	c.caches.Set(userID, user)

	return &user, nil
}

func (c *UserCollection) IsUser(ctx context.Context, userID int64) bool {
	user, err := c.FindByID(ctx, userID)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			slog.Error("error in gathering user is registered.", "error", err)
		}

		return false
	}

	c.caches.Set(userID, *user)

	return true
}

func (c *UserCollection) IsUserBlocked(ctx context.Context, userID int64) (bool, string) {
	user, err := c.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, ""
		}

		slog.Error("error in gathering user is blocked.", "error", err)
		return true, "에러가 발생하여 차단한 유저를 구별 못해요. 계속 이러면 연락주세요."
	}

	c.caches.Set(userID, *user)

	return user.Blocked, user.BlockedReason
}

func (c *UserCollection) GetUserChattingMode(ctx context.Context, userID int64) (ChattingMode, error) {
	user, err := c.FindByID(ctx, userID)
	if err != nil {
		slog.Error("error in gathering user's mode.", "error", err)
		return ChattingAIMode, err
	}

	c.caches.Set(userID, *user)

	return user.ChattingMode, nil
}

func (c *UserCollection) Update(ctx context.Context, userID int64, data *UserUpdate) (*mongo.UpdateResult, error) {
	result, err := c.coll.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return nil, err
	}

	c.caches.Delete(userID)

	// 캐시 저장용
	if _, err := c.FindByID(ctx, userID); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *UserCollection) Delete(ctx context.Context, userID int64) (*mongo.DeleteResult, error) {
	result, err := c.coll.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		return nil, err
	}

	c.caches.Delete(userID)

	return result, nil
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
