package repository

import (
	"context"
	"log"
	"time"

	"git.wh64.net/muffin/goMuffin/cache"
	"git.wh64.net/muffin/goMuffin/configs"
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

var instance *MuffinDatabase

func init() {
	const timeToExpire = time.Hour * 12

	client, err := mongo.Connect(options.Client().ApplyURI(configs.GetConfig().Database.URL))
	if err != nil {
		log.Panicln(err)
	}

	instance = &MuffinDatabase{
		Client:    client,
		Knowledge: &KnowledgeCollection{client.Database(configs.GetConfig().Database.Name).Collection("learn"), cache.New[*knowledgeCacheItem](timeToExpire)},
		Texts:     client.Database(configs.GetConfig().Database.Name).Collection("text"),
		Memory:    &MemoryCollection{client.Database(configs.GetConfig().Database.Name).Collection("memory"), cache.New[*memoryCacheItem](timeToExpire)},
		Users:     &UserCollection{client.Database(configs.GetConfig().Database.Name).Collection("user"), cache.New[*User](timeToExpire)},
		Chats:     &ChatCollection{client.Database(configs.GetConfig().Database.Name).Collection("chat")},
	}
}

func GetDatabase() *MuffinDatabase {
	return instance
}

func (d *MuffinDatabase) Disconnect() {
	GetDatabase().Client.Disconnect(context.TODO())
}
