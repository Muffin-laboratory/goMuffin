package chatbot

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
)

func (c *Chatbot) getMuffinResponse(ctx context.Context, question string) (string, error) {
	var result string
	x := rand.Intn(10)

	data, err := repository.GetDatabase().Texts.All(ctx)
	if err != nil {
		return "살려주ㅅ세요", err
	}

	knowledge, err := repository.GetDatabase().Knowledge.Find(ctx, query.KnowledgeQueryBuilder().SetCommand(question))
	if err != nil {
		return "살려주ㅅ세요", err
	}

	if x > 2 && len(knowledge) != 0 {
		data := knowledge[rand.Intn(len(knowledge))]
		user, _ := c.s.User(data.UserID)

		result =
			fmt.Sprintf("%s\n%s", data.Result, utils.InlineCode(fmt.Sprintf("%s님이 알려주셨어요.", user.Username)))
	} else {
		result = data[rand.Intn(len(data))].Text
	}
	return result, nil
}
