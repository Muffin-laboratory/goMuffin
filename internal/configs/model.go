package configs

import (
	"log/slog"

	"github.com/disgoorg/snowflake/v2"
)

type botConfig struct {
	Token   string
	Prefix  string
	OwnerID snowflake.ID
}

type trainConfig struct {
	UserID snowflake.ID
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
	Owner         string
	Repository    string
	OldRepository string
}

type integrateMDCConfig struct {
	Server struct {
		Port int
	}
}

type commandConfig struct {
	DeveloperOnlyGuildID snowflake.ID
}

type loggerConfig struct {
	Level     slog.Level
	WriteFile bool
}
