package migration

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type originalUser struct {
	ID                        bson.ObjectID           `bson:"_id,omitempty"`
	UserID                    string                  `bson:"user_id,omitempty"`
	Blocked                   bool                    `bson:"blocked"`
	BlockedReason             string                  `bson:"blocked_reason"`
	ChatID                    bson.ObjectID           `bson:"chat_id"`
	CreatedAt                 time.Time               `bson:"created_at"`
	ChattingMode              repository.ChattingMode `bson:"chatting_mode"`
	ReplyUser                 bool                    `bson:"reply_user"`
	CreateNewChatAfter12Hours bool                    `bson:"create_new_chat_after_12_hours"`
}

func migrateUser(coll *mongo.Collection, ch chan *migrationErr, wg *sync.WaitGroup) {
	const where = "user"

	defer wg.Done()

	var dataToMigrate []originalUser
	var migrationData []repository.User
	var idsToDelete []bson.D

	cur, err := coll.Find(context.Background(), bson.D{
		{
			Key:   "_id",
			Value: bson.M{"$type": "objectId"},
		},
	})
	if err != nil {
		ch <- &migrationErr{where, err}
		return
	}

	defer cur.Close(context.Background())

	if err = cur.All(context.Background(), &dataToMigrate); err != nil {
		ch <- &migrationErr{where, err}
		return
	}

	if len(dataToMigrate) == 0 {
		ch <- nil
		return
	}

	for _, data := range dataToMigrate {
		id, err := strconv.ParseInt(data.UserID, 10, 0)
		if err != nil {
			ch <- &migrationErr{where, err}
			return
		}

		idsToDelete = append(idsToDelete, bson.D{{Key: "_id", Value: data.ID}})
		migrationData = append(migrationData, repository.User{
			ID:                        id,
			Blocked:                   data.Blocked,
			BlockedReason:             data.BlockedReason,
			ChatID:                    data.ChatID,
			CreatedAt:                 data.CreatedAt,
			ChattingMode:              data.ChattingMode,
			ReplyUser:                 data.ReplyUser,
			CreateNewChatAfter12Hours: data.CreateNewChatAfter12Hours,
		})
	}

	_, err = coll.InsertMany(context.Background(), migrationData)
	if err != nil {
		ch <- &migrationErr{where, err}
		return
	}

	for _, filter := range idsToDelete {
		_, err = coll.DeleteMany(context.Background(), filter)
		if err != nil {
			ch <- &migrationErr{where, err}
			return
		}
	}

	ch <- nil
}
