package configs

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/disgoorg/snowflake/v2"
)

func getRequiredValue(key string) string {
	value := os.Getenv(key)
	if value == "" {
		slog.Error("[Fatal] a required env value is not found.", "key", key)
		os.Exit(1)
	}

	return value
}

func getRequiredValueToSnowflake(key string) snowflake.ID {
	value := getRequiredValue(key)
	id, err := snowflake.Parse(value)
	if err != nil {
		slog.Error("[Fatal] failed to get a required snowflake env value.", "key", key, "value", value)
		os.Exit(1)
	}

	return id
}

func getValue(key string) string {
	return os.Getenv(key)
}

func getValueToSnowflake(key string) snowflake.ID {
	value := getValue(key)
	if value == "" {
		return 0
	}

	id, err := snowflake.Parse(value)
	if err != nil {
		slog.Error("[Fatal] failed to get a snowflake env value.", "key", key, "value", value)
		os.Exit(1)
	}

	return id
}

func getValueToInt(key string) int {
	value := getValue(key)
	if value == "" {
		return 0
	}

	parsedInt, err := strconv.Atoi(value)
	if err != nil {
		slog.Error("[Fatal] failed to get an int env value.", "key", key, "value", value)
		os.Exit(1)
	}

	return parsedInt
}

func getValueToBool(key string) bool {
	value := getValue(key)
	if value == "" {
		return false
	}

	parsedBool, err := strconv.ParseBool(value)
	if err != nil {
		slog.Error("[Fatal] failed to get a bool env value.", "key", key, "value", value)
		os.Exit(1)
	}

	return parsedBool
}

func getValueToLogLevel(key string) slog.Level {
	value := getValue(key)

	switch value {
	case "ERROR":
		return slog.LevelError
	case "WARN":
		return slog.LevelWarn
	case "DEBUG":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}
