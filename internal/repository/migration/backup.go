package migration

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func BackupDB(db *mongo.Database) {
	var errored bool
	const collections = 5

	var wg sync.WaitGroup

	slog.Info("start creating database backup...")

	ch := make(chan *migrationErr, collections)
	wg.Add(collections)

	go backup("text", db.Collection("text"), ch, &wg)
	go backup("knowledge", db.Collection("learn"), ch, &wg)
	go backup("memory", db.Collection("memory"), ch, &wg)
	go backup("user", db.Collection("user"), ch, &wg)
	go backup("chat", db.Collection("chat"), ch, &wg)

	wg.Wait()
	close(ch)
	for err := range ch {
		if err != nil {
			slog.Error("error while creating backup", "where", err.where, "error", err)
			errored = true
		}
	}

	if errored {
		os.Exit(1)
	}

	slog.Info("database backup is success.")
}

func backup(where string, coll *mongo.Collection, ch chan *migrationErr, wg *sync.WaitGroup) {
	defer wg.Done()

	now := time.Now().Format("060102")

	cur, err := coll.Find(context.Background(), bson.D{})
	if err != nil {
		ch <- &migrationErr{where, err}
		return
	}

	defer func() {
		if err := cur.Close(context.Background()); err != nil {
			slog.Error("failed to close "+where+" cursor while backup database", "error", err)
		}
	}()

	var data []map[string]any
	if err = cur.All(context.Background(), &data); err != nil {
		ch <- &migrationErr{where, err}
		return
	}

	content, err := json.Marshal(data)
	if err != nil {
		ch <- &migrationErr{where, err}
		return
	}

	if err = os.WriteFile(fmt.Sprintf("backup/%s_%s.json", where, now), content, 0777); err != nil {
		ch <- &migrationErr{where, err}
		return
	}

	ch <- nil
}
