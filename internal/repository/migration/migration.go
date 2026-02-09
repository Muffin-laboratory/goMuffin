package migration

import (
	"log/slog"
	"sync"

	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type migrationErr struct {
	where string
	err   error
}

func MigrationUserIDToInt64(client *mongo.Client) {
	const collections = 4

	var wg sync.WaitGroup

	db := client.Database(configs.GetConfig().Database.Name)
	ch := make(chan *migrationErr, collections)
	wg.Add(collections)

	go migrateUser(db.Collection("user"), ch, &wg)
	go migrateAnotherCollection(db.Collection("memory"), ch, &wg)
	go migrateAnotherCollection(db.Collection("chat"), ch, &wg)
	go migrateAnotherCollection(db.Collection("learn"), ch, &wg)

	wg.Wait()
	close(ch)
	for err := range ch {
		if err != nil {
			slog.Error("error while migration.", "where", err.where, "error", err)
		}
	}

	slog.Info("migration is success!!!")
}
