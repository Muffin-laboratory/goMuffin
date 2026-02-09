package main

import (
	"log/slog"
	"os"

	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/migration"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	client, err := mongo.Connect(options.Client().ApplyURI(configs.GetConfig().Database.URL))
	if err != nil {
		slog.Error("error while create database.", "error", err)
		os.Exit(1)
	}

	migration.MigrationUserIDToInt64(client)
}
