package chatbot

import (
	"context"
	"fmt"
	"os"

	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/disgoorg/disgo/discord"
)

func loadPrompt() (string, error) {
	bin, err := os.ReadFile(configs.GetConfig().Chatbot.Gemini.PromptPath)
	if err != nil {
		return "", err
	}

	return string(bin), nil
}

func makePrompt(ctx context.Context, systemPrompt string, user *discord.User) (string, error) {
	var userPrompt string

	knowledgePrompt := "## Knowledge of the user\n"

	if user.ID == configs.GetConfig().Bot.OwnerID {
		userPrompt += fmt.Sprintf(
			"---\n## User Information\n* **ID:** %s\n* **Name:** %s\n* **Other:** This user is your developer.",
			user.ID.String(),
			*user.GlobalName,
		)
	} else {
		userPrompt += fmt.Sprintf(
			"---\n## User Information\n* **ID:** %s\n* **Name:** %s\n* **Other:** This user is **not** your developer.",
			user.ID.String(),
			*user.GlobalName,
		)
	}

	knowledge, err := repository.GetDatabase().Knowledge.Find(ctx, query.KnowledgeQueryBuilder().SetUserID(user.ID.String()))
	if err != nil {
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
