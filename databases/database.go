package databases

import (
	"context"
	"log"

	"git.wh64.net/muffin/goMuffin/configs"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MuffinDatabase struct {
	Client *mongo.Client
	Learns *mongo.Collection
	Texts  *mongo.Collection
	Memory *mongo.Collection
	Users  *mongo.Collection
	Chats  *mongo.Collection
}

var instance *MuffinDatabase

func init() {
	client, err := mongo.Connect(options.Client().ApplyURI(configs.GetConfig().Database.URL))
	if err != nil {
		log.Panicln(err)
	}
	instance = &MuffinDatabase{
		Client: client,
		Learns: client.Database(configs.GetConfig().Database.Name).Collection("learn"),
		Texts:  client.Database(configs.GetConfig().Database.Name).Collection("text"),
		Memory: client.Database(configs.GetConfig().Database.Name).Collection("memory"),
		Users:  client.Database(configs.GetConfig().Database.Name).Collection("user"),
		Chats:  client.Database(configs.GetConfig().Database.Name).Collection("chat"),
	}
}

func GetDatabase() *MuffinDatabase {
	return instance
}

func Disconnect() {
	GetDatabase().Client.Disconnect(context.TODO())
}
