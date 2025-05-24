package databases

import (
	"log"

	"git.wh64.net/muffin/goMuffin/configs"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MuffinDatabase struct {
	Client *mongo.Client
	Learns *mongo.Collection
	Texts  *mongo.Collection
}

var Database *MuffinDatabase

func init() {
	var err error

	Database, err = Connect()
	if err != nil {
		log.Fatalln(err)
	}
}

func Connect() (*MuffinDatabase, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(configs.Config.DatabaseURL))
	if err != nil {
		return nil, err
	}
	return &MuffinDatabase{
		Client: client,
		Learns: client.Database(configs.Config.DatabaseName).Collection("learn"),
		Texts:  client.Database(configs.Config.DatabaseName).Collection("text"),
	}, nil
}
