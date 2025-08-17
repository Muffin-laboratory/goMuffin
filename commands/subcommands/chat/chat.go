package chat

import (
	"log"

	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/utils"
)

func Chat(i *utils.InteractionCreate, content string) error {
	str, err := chatbot.GetChatBot().GetResponse(i.User, content)
	if err != nil {
		log.Println(err)
		i.EditReply(&utils.InteractionEdit{
			Content: &str,
		})
		return nil
	}

	result := chatbot.ParseResult(str, i.Session, i)
	return i.EditReply(&utils.InteractionEdit{
		Content: &result,
	})
}
