package configs

import "fmt"

func AddPrefix(command string) string {
	return fmt.Sprintf("%s%s", GetConfig().Bot.Prefix, command)
}
