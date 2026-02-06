package chat

import (
	"log"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/chatbot"
	"github.com/bwmarrin/discordgo"
)

func Chat(i *builders.InteractionCreate, opts *discordgo.ApplicationCommandInteractionDataOption) error {
	var attachment *discordgo.MessageAttachment

	if opt := opts.GetOption("첨부파일"); opt != nil {
		id := opt.Value.(string)

		if resolvedData := i.ApplicationCommandData().Resolved; resolvedData != nil && resolvedData.Attachments != nil {
			attachment = resolvedData.Attachments[id]
		}
	}

	content, err := chatbot.GetChatBot().GetResponse(i.Ctx, i.User, opts.GetOption("내용").StringValue(), attachment)
	if err != nil {
		log.Println(err)
		i.EditReply(&discordgo.WebhookEdit{
			Content: &content,
		})
		return nil
	}

	result := chatbot.ParseResult(content, i.Session, i)
	return i.EditReply(&discordgo.WebhookEdit{
		Content: &result,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{},
			Users: []string{},
			Roles: []string{},
		},
	})
}
