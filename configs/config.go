package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

type databaseConfig struct {
	Name       string
	URL        string
	HostName   string
	Username   string
	Password   string
	AuthSource string
	Port       int
}

// MuffinConfig for Muffin bot
type MuffinConfig struct {
	Bot      botConfig
	Train    trainConfig
	Database databaseConfig

	// Deprecated: Use Database.URL
	DatabaseURL string

	// Deprecated: Use Database.Name
	DatabaseName string
}

var Config *MuffinConfig

func init() {
	godotenv.Load()
	Config = &MuffinConfig{Bot: botConfig{}, Train: trainConfig{}, Database: databaseConfig{}}
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

	config.Database.URL = getValue("DATABASE_URL")
	config.Database.HostName = getValue("DATABASE_HOSTNAME")
	config.Database.Password = getValue("DATABASE_PASSWORD")
	config.Database.Username = getValue("DATABASE_USERNAME")
	config.Database.AuthSource = getValue("DATABASE_AUTH_SOURCE")
	config.Database.Name = getRequiredValue("DATABASE_NAME")
	port, err := strconv.Atoi(getValue("DATABASE_PORT"))
	if err != nil {
		log.Println("[goMuffin] 'DATABASE_PORT'값을 int로 파싱할 수 없어요.")
		log.Fatalln(err)
	}

	config.Database.Port = port

	if config.Database.AuthSource == "" {
		config.Database.AuthSource = "admin"
	}

	if config.Database.URL == "" {
		config.Database.URL = fmt.Sprintf("mongodb://%s:%s@%s:%d/?authSource=%s", config.Database.Username, config.Database.Password, config.Database.HostName, config.Database.Port, config.Database.AuthSource)
	}
	fmt.Println(config.Database.URL)

	// Deprecated된 Value
	config.DatabaseURL = config.Database.URL
	config.DatabaseName = config.Database.Name
}
