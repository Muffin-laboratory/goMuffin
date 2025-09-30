package chat

import (
	"log"

	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func Chat(i *utils.InteractionCreate, opts utils.CommandInteractionOptionsMap) error {
	str, err := chatbot.GetChatBot().GetResponse(i.User, opts["내용"].StringValue())
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
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{},
			Users: []string{},
			Roles: []string{},
		},
	})
}
