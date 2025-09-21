package chatbot

import (
	"context"
	"fmt"
	"os"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"github.com/bwmarrin/discordgo"
)

func loadPrompt() (string, error) {
	bin, err := os.ReadFile(configs.GetConfig().Chatbot.Gemini.PromptPath)
	if err != nil {
		return "", err
	}

	return string(bin), nil
}

func makePrompt(systemPrompt string, user *discordgo.User) (string, error) {
	var knowledge []databases.Knowledge
	var userPrompt string

	knowledgePrompt := "## Knowledge of the user\n"

	if user.ID == configs.GetConfig().Bot.OwnerID {
		userPrompt += fmt.Sprintf(
			"---\n## User Information\n* **ID:** %s\n* **Name:** %s\n* **Other:** This user is your developer.",
			user.ID,
			user.GlobalName,
		)
	} else {
		userPrompt += fmt.Sprintf(
			"---\n## User Information\n* **ID:** %s\n* **Name:** %s\n* **Other:** This user is **not** your developer.",
			user.ID,
			user.GlobalName,
		)
	}

	cur, err := databases.GetDatabase().Knowledge.Find(context.TODO(), databases.Knowledge{UserID: user.ID})
	if err != nil {
		return "", err
	}

	if err = cur.All(context.TODO(), &knowledge); err != nil {
		return "", err
	}

	if len(knowledge) == 0 {
		knowledgePrompt += "* **Knowledge of the user(the user is not Muffin.) is None.**"
	}

	for _, knowledge := range knowledge {
		knowledgePrompt += fmt.Sprintf("* **%s**: %s\n", knowledge.Command, knowledge.Result)
	}

	return systemPrompt + userPrompt + knowledgePrompt, nil
}
