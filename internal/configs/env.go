package configs

import (
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/disgoorg/snowflake/v2"
)

func getRequiredValue(key string) string {
	value := os.Getenv(key)
	if value == "" {
		slog.Error("[Fatal] an required env value is not found.", "key", key)
		os.Exit(1)
	}

	return value
}

func getRequiredValueToSnowflake(key string) snowflake.ID {
	value := getRequiredValue(key)
	id, err := snowflake.Parse(value)
	if err != nil {
		slog.Error("[Fatal] failed to required snowflake env value.", "key", key, "value", value)
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
		slog.Error("[Fatal] failed to snowflake env value.", "key", key, "value", value)
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
		log.Fatalf("[goMuffin] .env 파일에서 '%s'값은 int타입이어야 해요.", key)
	}

	return parsedInt
}
