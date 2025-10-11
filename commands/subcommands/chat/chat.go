package chat

import (
	"log"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/chatbot"
	"github.com/bwmarrin/discordgo"
)

func Chat(i *builders.InteractionCreate, opts builders.CommandInteractionOptionsMap) error {
	str, err := chatbot.GetChatBot().GetResponse(i.User, opts["내용"].StringValue())
	if err != nil {
		log.Println(err)
		i.EditReply(&builders.InteractionEdit{
			Content: &str,
		})
		return nil
	}

	result := chatbot.ParseResult(str, i.Session, i)
	return i.EditReply(&builders.InteractionEdit{
		Content: &result,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{},
			Users: []string{},
			Roles: []string{},
		},
	})
}
