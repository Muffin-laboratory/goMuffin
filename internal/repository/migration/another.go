package migration

import (
	"context"
	"strconv"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func migrateAnotherCollection(coll *mongo.Collection, ch chan *migrationErr, wg *sync.WaitGroup) {
	const where = "memory"

	defer wg.Done()

	var dataToMigrate []map[string]any
	migrationData := make(map[any]bson.D)

	cur, err := coll.Find(context.Background(), bson.D{
		{
			Key:   "user_id",
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
		id, err := strconv.ParseInt(data["user_id"].(string), 10, 0)
		if err != nil {
			ch <- &migrationErr{where, err}
			return
		}

		migrationData[data["_id"]] = bson.D{{
			Key:   "$set",
			Value: bson.M{"user_id": id},
		}}
	}

	for key, data := range migrationData {
		_, err := coll.UpdateByID(context.Background(), key, data)
		if err != nil {
			ch <- &migrationErr{where, err}
			return
		}
	}

	ch <- nil
}
