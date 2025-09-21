package databases

import (
	"context"
	"log"

	"git.wh64.net/muffin/goMuffin/configs"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MuffinDatabase struct {
	Client    *mongo.Client
	Knowledge *mongo.Collection
	Texts     *mongo.Collection
	Memory    *MemoryCollection
	Users     *UserCollection
	Chats     *ChatCollection
}

var instance *MuffinDatabase

func init() {
	client, err := mongo.Connect(options.Client().ApplyURI(configs.GetConfig().Database.URL))
	if err != nil {
		log.Panicln(err)
	}

	instance = &MuffinDatabase{
		Client:    client,
		Knowledge: client.Database(configs.GetConfig().Database.Name).Collection("learn"),
		Texts:     client.Database(configs.GetConfig().Database.Name).Collection("text"),
		Memory:    &MemoryCollection{client.Database(configs.GetConfig().Database.Name).Collection("memory")},
		Users:     &UserCollection{client.Database(configs.GetConfig().Database.Name).Collection("user")},
		Chats:     &ChatCollection{client.Database(configs.GetConfig().Database.Name).Collection("chat")},
	}
}

func GetDatabase() *MuffinDatabase {
	return instance
}

func (d *MuffinDatabase) Disconnect() {
	GetDatabase().Client.Disconnect(context.TODO())
}
