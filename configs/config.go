package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type botConfig struct {
	Token   string
	Prefix  string
	OwnerId string
}

type trainConfig struct {
	UserID string
}

// MuffinConfig for Muffin bot
type MuffinConfig struct {
	Bot          botConfig
	Train        trainConfig
	DatabaseURL  string
	DatabaseName string
}

var Config *MuffinConfig

func init() {
	godotenv.Load()
	Config = &MuffinConfig{Bot: botConfig{}, Train: trainConfig{}}
	setConfig(Config)
}

func getRequiredValue(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalln(fmt.Sprintf("[goMuffin] .env 파일에서 필요한 '%s'값이 없어요.", key))
	}
	return value
}

func getValue(key string) string {
	return os.Getenv(key)
}

func setConfig(config *MuffinConfig) {
	config.Bot.Prefix = getRequiredValue("BOT_PREFIX")
	config.Bot.Token = getRequiredValue("BOT_TOKEN")
	config.Bot.OwnerId = getRequiredValue("BOT_OWNER_ID")

	config.Train.UserID = getValue("TRAIN_USER_ID")

	config.DatabaseURL = getRequiredValue("DATABASE_URL")
	config.DatabaseName = getRequiredValue("DATABASE_NAME")
}
