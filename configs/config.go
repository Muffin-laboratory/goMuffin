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
	OwnerID string
}

type trainConfig struct {
	UserID string
}

type geminiConfig struct {
	Token      string
	PromptPath string
	Model      string
}

type chatbotConfig struct {
	Train  trainConfig
	Gemini geminiConfig
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

type serviceConfig struct {
	PrivacyPolicyURL string
	TermOfServiceURL string
}

type githubConfig struct {
	Owner      string
	Repository string
}

// MuffinConfig for Muffin bot
type MuffinConfig struct {
	Bot      botConfig
	Database databaseConfig
	Chatbot  chatbotConfig
	Service  serviceConfig
	GitHub   githubConfig
}

var instance *MuffinConfig

func init() {
	godotenv.Load()
	instance = &MuffinConfig{}
	setConfig(instance)
}

func GetConfig() *MuffinConfig {
	return instance
}

func getRequiredValue(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("[goMuffin] .env 파일에서 필요한 '%s'값이 없어요.", key)
	}
	return value
}

func getValue(key string) string {
	return os.Getenv(key)
}

func setConfig(config *MuffinConfig) {
	config.Bot = botConfig{
		Prefix:  getRequiredValue("BOT_PREFIX"),
		Token:   getRequiredValue("BOT_TOKEN"),
		OwnerID: getRequiredValue("BOT_OWNER_ID"),
	}

	config.Database = databaseConfig{
		URL:        getValue("DATABASE_URL"),
		HostName:   getValue("DATABASE_HOSTNAME"),
		Password:   getValue("DATABASE_PASSWORD"),
		Username:   getValue("DATABASE_USERNAME"),
		AuthSource: getValue("DATABASE_AUTH_SOURCE"),
		Name:       getRequiredValue("DATABASE_NAME"),
	}
	port, err := strconv.Atoi(getValue("DATABASE_PORT"))
	if getValue("DATABASE_PORT") != "" && err != nil {
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

	config.Chatbot = chatbotConfig{
		Gemini: geminiConfig{Token: getValue("CHATBOT_GEMINI_TOKEN"), PromptPath: getValue("CHATBOT_GEMINI_PROMPT_PATH"), Model: getValue("CHATBOT_GEMINI_MODEL")},
		Train:  trainConfig{UserID: getValue("CHATBOT_TRAIN_USER_ID")},
	}

	if config.Chatbot.Gemini.Model == "" {
		config.Chatbot.Gemini.Model = "gemini-2.0-flash"
	}

	config.Service = serviceConfig{
		PrivacyPolicyURL: getRequiredValue("SERVICE_PRIVACY_POLICY_URL"),
		TermOfServiceURL: getRequiredValue("SERVICE_TERM_OF_SERVICE_URL"),
	}

	config.GitHub = githubConfig{
		Owner:      getValue("GITHUB_OWNER"),
		Repository: getValue("GITHUB_REPO"),
	}
}
