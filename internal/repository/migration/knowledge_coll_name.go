package migration

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func MigrateKnowledgeCollectionName(db *mongo.Database) {
	ch := make(chan error, 1)

	slog.Info("start knowledge collection migration...")

	go func() {
		const where = "knowledge_name"

		var data []map[string]any

		oldColl := db.Collection("learn")
		newColl := db.Collection("knowledge")

		cur, err := oldColl.Find(context.Background(), bson.D{})
		if err != nil {
			ch <- err
			return
		}

		defer cur.Close(context.Background())

		if err = cur.All(context.Background(), &data); err != nil {
			ch <- err
			return
		}

		result, err := newColl.InsertMany(context.Background(), data)
		if err != nil {
			ch <- err
			return
		}

		if len(result.InsertedIDs) == 0 {
			ch <- fmt.Errorf("failed to insert data in new collection")
			return
		}

		if err = oldColl.Drop(context.Background()); err != nil {
			ch <- err
			return
		}

		ch <- nil
	}()

	if err := <-ch; err != nil {
		slog.Error("error while change knowledge collection name.", "error", err)
		os.Exit(1)
	}

	slog.Info("knowledge collection name migration is success!!!")
}
