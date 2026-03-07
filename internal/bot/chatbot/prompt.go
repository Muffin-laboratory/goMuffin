package chatbot

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/disgoorg/disgo/discord"
)

func loadPrompt() (string, string, error) {
	defaultPrompt, err := os.ReadFile(configs.Configs().Chatbot.Gemini.PromptPath)
	if err != nil {
		return "", "", err
	}

	corePrompt, err := os.ReadFile("core_prompt.txt")
	if err != nil {
		return "", "", err
	}

	return string(defaultPrompt), string(corePrompt), nil
}

func makePrompt(ctx context.Context, systemPrompt, corePrompt string, user *discord.User, isChat bool) (string, error) {
	var userPrompt string
	var knowledgePrompt strings.Builder

	if user.ID == configs.Configs().Bot.OwnerID {
		userPrompt += fmt.Sprintf(
			"\n---\n## User Information\n* **ID:** %s\n* **Name:** %s\n* **Other:** This user is your developer.",
			user.ID.String(),
			*user.GlobalName,
		)
	} else {
		userPrompt += fmt.Sprintf(
			"\n---\n## User Information\n* **ID:** %s\n* **Name:** %s\n* **Other:** This user is **not** your developer.",
			user.ID.String(),
			*user.GlobalName,
		)
	}

	if isChat {
		knowledgePrompt.WriteString("## Knowledge of the user\n")

		knowledge, err := repository.GetDatabase().Knowledge.Find(ctx, query.KnowledgeQueryBuilder().SetUserID(int64(user.ID)))
		if err != nil {
			return "", err
		}

		if len(knowledge) == 0 {
			knowledgePrompt.WriteString("* **Knowledge of the user(The user is not You.) is None.**")
		}

		for _, knowledge := range knowledge {
			fmt.Fprintf(&knowledgePrompt, "* **%s**: %s\n", knowledge.Command, knowledge.Result)
		}
	}

	return systemPrompt + userPrompt + knowledgePrompt.String() + corePrompt, nil
}
