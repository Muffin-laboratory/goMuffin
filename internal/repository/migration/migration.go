package migration

import (
	"log/slog"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type migrationErr struct {
	where string
	err   error
}

func MigrationUserIDToInt64(db *mongo.Database) {
	var errored bool
	const collections = 4

	var wg sync.WaitGroup

	slog.Info("start user id migration...")

	ch := make(chan *migrationErr, collections)
	wg.Add(collections)

	go migrateUser(db.Collection("user"), ch, &wg)
	go migrateAnotherCollection("memory", db.Collection("memory"), ch, &wg)
	go migrateAnotherCollection("chat", db.Collection("chat"), ch, &wg)
	go migrateAnotherCollection("knowledge", db.Collection("learn"), ch, &wg)

	wg.Wait()
	close(ch)
	for err := range ch {
		if err != nil {
			slog.Error("error while user id migration.", "where", err.where, "error", err)
			errored = true
		}
	}

	if errored {
		os.Exit(1)
	}

	slog.Info("user id migration is success!!!")
}
