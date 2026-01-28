package chat

import (
	"log"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/chatbot"
	"github.com/bwmarrin/discordgo"
)

func Chat(i *builders.InteractionCreate, opts builders.CommandInteractionOptionsMap) error {
	var attachment *discordgo.MessageAttachment

	if opt, ok := opts["첨부파일"]; ok {
		id := opt.Value.(string)

		if resolvedData := i.ApplicationCommandData().Resolved; resolvedData != nil && resolvedData.Attachments != nil {
			attachment = resolvedData.Attachments[id]
		}
	}

	str, err := chatbot.GetChatBot().GetResponse(i.Ctx, i.User, opts["내용"].StringValue(), attachment)
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
