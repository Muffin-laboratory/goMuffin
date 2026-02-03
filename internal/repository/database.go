package repository

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/cache"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MuffinDatabase struct {
	Client    *mongo.Client
	Knowledge *KnowledgeCollection
	Texts     *mongo.Collection
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
			log.Panicln(err)
		}

		instance = &MuffinDatabase{
			Client:    client,
			Knowledge: &KnowledgeCollection{client.Database(configs.GetConfig().Database.Name).Collection("learn"), cache.New[string, *knowledgeCacheItem](timeToExpire)},
			Texts:     client.Database(configs.GetConfig().Database.Name).Collection("text"),
			Memory:    &MemoryCollection{client.Database(configs.GetConfig().Database.Name).Collection("memory"), cache.New[string, *memoryCacheItem](timeToExpire)},
			Users:     &UserCollection{client.Database(configs.GetConfig().Database.Name).Collection("user"), cache.New[string, User](timeToExpire)},
			Chats:     newChatCollection(client.Database(configs.GetConfig().Database.Name).Collection("chat")),
		}
	})
	return instance
}

func (d *MuffinDatabase) Disconnect() {
	GetDatabase().Client.Disconnect(context.TODO())
}
