package chatbot

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (c *Chatbot) getMuffinResponse(ctx context.Context, question string) (string, error) {
	var data []repository.Text
	var result string
	x := rand.Intn(10)

	cur, err := repository.GetDatabase().Texts.Find(context.TODO(), bson.D{{Key: "persona", Value: "muffin"}})
	if err != nil {
		return "살려주ㅅ세요", err
	}

	defer cur.Close(context.TODO())

	if err = cur.All(context.TODO(), &data); err != nil {
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
