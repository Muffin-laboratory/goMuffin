package repository

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/cache"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MuffinDatabase struct {
	Client    *mongo.Client
	Knowledge *KnowledgeCollection
	Texts     *TextCollection
	Memory    *MemoryCollection
	Users     *UserCollection
	Chats     *ChatCollection
}

const timeToExpire = time.Hour * 12

var instance *MuffinDatabase
var once sync.Once

func GetDatabase() *MuffinDatabase {
	once.Do(func() {
		client, err := mongo.Connect(options.Client().ApplyURI(configs.GetConfig().Database.URL))
		if err != nil {
			slog.Error("error while create database.", "error", err)
			os.Exit(1)
		}

		instance = &MuffinDatabase{
			Client:    client,
			Knowledge: newKnowledgeCollection(client.Database(configs.GetConfig().Database.Name).Collection("knowledge")),
			Texts:     newTextCollection(client.Database(configs.GetConfig().Database.Name).Collection("text")),
			Memory:    newMemoryCollection(client.Database(configs.GetConfig().Database.Name).Collection("memory")),
			Users:     newUserCollection(client.Database(configs.GetConfig().Database.Name).Collection("user")),
			Chats:     newChatCollection(client.Database(configs.GetConfig().Database.Name).Collection("chat")),
		}

		slog.Info("database is created.")
	})
	return instance
}

func (d *MuffinDatabase) Disconnect() {
	GetDatabase().Client.Disconnect(context.Background())
}

func createIndexCache(im *cache.CacheManager[string, *indexItem], id bson.ObjectID, index indexBuilder) {
	if cache, ok := im.Get(index.build()); ok {
		cache.mu.Lock()
		cache.ids = append(cache.ids, id)
		cache.mu.Unlock()
	} else {
		im.Set(index.build(), indexItemBuilder(id))
	}
}
