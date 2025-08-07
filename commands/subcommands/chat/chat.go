package chat

import (
	"log"

	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func Chat(m any, user *discordgo.User, content string) error {
	// 채팅하기는 슬래시 커맨드만 가능
	i := m.(*utils.InteractionCreate)

	str, err := chatbot.ChatBot.GetResponse(user, content)
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
