package configs

import (
	"fmt"
	"sync"

	"github.com/joho/godotenv"
)

// MuffinConfig for Muffin bot
type MuffinConfig struct {
	Bot          botConfig
	Database     databaseConfig
	Chatbot      chatbotConfig
	Service      serviceConfig
	GitHub       githubConfig
	IntegrateMDC integrateMDCConfig
	Command      commandConfig
	Logger       loggerConfig
}

var instance *MuffinConfig
var once sync.Once

func Configs() *MuffinConfig {
	once.Do(func() {
		_ = godotenv.Load()
		instance = &MuffinConfig{}
		setConfig(instance)
	})

	return instance
}

func setConfig(config *MuffinConfig) {
	config.Bot = botConfig{
		Prefix:  getRequiredValue("BOT_PREFIX"),
		Token:   getRequiredValue("BOT_TOKEN"),
		OwnerID: getRequiredValueToSnowflake("BOT_OWNER_ID"),
	}

	config.Database = databaseConfig{
		HostName:   getRequiredValue("DATABASE_HOSTNAME"),
		Password:   getRequiredValue("DATABASE_PASSWORD"),
		Username:   getRequiredValue("DATABASE_USERNAME"),
		AuthSource: getRequiredValue("DATABASE_AUTH_SOURCE"),
		Name:       getRequiredValue("DATABASE_NAME"),
		Port:       getValueToInt("DATABASE_PORT"),
	}

	if config.Database.AuthSource == "" {
		config.Database.AuthSource = "admin"
	}

	if config.Database.Port == 0 {
		config.Database.Port = 21017
	}

	config.Database.URL = fmt.Sprintf("mongodb://%s:%s@%s:%d/?authSource=%s", config.Database.Username, config.Database.Password, config.Database.HostName, config.Database.Port, config.Database.AuthSource)

	config.Chatbot = chatbotConfig{
		Gemini: geminiConfig{
			Token:      getValue("CHATBOT_GEMINI_TOKEN"),
			PromptPath: getValue("CHATBOT_GEMINI_PROMPT_PATH"),
			Model:      getValue("CHATBOT_GEMINI_MODEL"),
		},
		Train: trainConfig{UserID: getValueToSnowflake("CHATBOT_TRAIN_USER_ID")},
	}

	if config.Chatbot.Gemini.Model == "" {
		config.Chatbot.Gemini.Model = "gemini-3-flash-preview"
	}

	config.Service = serviceConfig{
		PrivacyPolicyURL: getRequiredValue("SERVICE_PRIVACY_POLICY_URL"),
		TermOfServiceURL: getRequiredValue("SERVICE_TERM_OF_SERVICE_URL"),
	}

	config.GitHub = githubConfig{
		Owner:         getValue("GITHUB_OWNER"),
		Repository:    getValue("GITHUB_REPO"),
		OldRepository: getValue("GITHUB_OLD_REPO"),
	}

	config.IntegrateMDC = integrateMDCConfig{
		Server: struct{ Port int }{
			Port: getValueToInt("INTEGRATE_MDC_SERVER_PORT"),
		},
	}

	config.Command = commandConfig{
		DeveloperOnlyGuildID: getRequiredValueToSnowflake("COMMAND_DEVELOPER_ONLY_GUILD_ID"),
	}

	config.Logger = loggerConfig{
		Level:     getValueToLogLevel("LOGGER_LEVEL"),
		WriteFile: getValueToBool("LOGGER_WRITE_FILE"),
	}
}
