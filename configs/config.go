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
	Bot         botConfig
	Train       trainConfig
	DatabaseURL string
	DBName      string
}

func loadConfig() *MuffinConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("[goMuffin] 봇의 설절파일을 불러올 수가 없어요.")
		log.Fatalln(err)
	}
	config := &MuffinConfig{Bot: botConfig{}, Train: trainConfig{}}
	setConfig(config)

	return config
}

func getRequiredValue(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalln(fmt.Sprintf("[goMuffin] .env 파일에서 필요한 %s값이 없어요.", key))
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
	config.DBName = getRequiredValue("DATABASE_NAME")
}

var Config *MuffinConfig = loadConfig()
