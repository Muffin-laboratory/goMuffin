package repository

import (
	"context"
	"log"
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
			log.Panicln(err)
		}

		instance = &MuffinDatabase{
			Client:    client,
			Knowledge: newKnowledgeCollection(client.Database(configs.GetConfig().Database.Name).Collection("learn")),
			Texts:     newTextCollection(client.Database(configs.GetConfig().Database.Name).Collection("text")),
			Memory:    newMemoryCollection(client.Database(configs.GetConfig().Database.Name).Collection("memory")),
			Users:     &UserCollection{client.Database(configs.GetConfig().Database.Name).Collection("user"), cache.New[string, User](timeToExpire)},
			Chats:     newChatCollection(client.Database(configs.GetConfig().Database.Name).Collection("chat")),
		}
	})
	return instance
}

func (d *MuffinDatabase) Disconnect() {
	GetDatabase().Client.Disconnect(context.TODO())
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
