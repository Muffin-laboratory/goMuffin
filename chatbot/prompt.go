package chatbot

import (
	"fmt"
	"os"

	"git.wh64.net/muffin/goMuffin/configs"
	"github.com/bwmarrin/discordgo"
)

func loadPrompt() (string, error) {
	bin, err := os.ReadFile(configs.Config.Chatbot.Gemini.PromptPath)
	if err != nil {
		return "", err
	}

	return string(bin), nil
}

func makePrompt(systemPrompt string, user *discordgo.User) string {
	if user.ID == configs.Config.Bot.OwnerId {
		return fmt.Sprintf(systemPrompt, fmt.Sprintf(
			"## User Information\n* **ID:** %s\n* **Name:** %s\n* **Other:** This user is your developer.",
			user.ID,
			user.GlobalName,
		))
	}

	return fmt.Sprintf(systemPrompt, fmt.Sprintf(
		"## User Information\n* **ID:** %s\n* **Name:** %s\n* **Other:** This user is **not** your developer.",
		user.ID,
		user.GlobalName,
	))
}
