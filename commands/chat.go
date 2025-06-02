package commands

import (
	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var ChatCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "대화",
		Description: "이 봇이랑 대화해요. (슬래시 커맨드 전용)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "내용",
				Description: "대화할 내용",
				Required:    true,
			},
		},
	},
	DetailedDescription: &DetailedDescription{
		Usage:    "/대화 (내용:대화할 내용)",
		Examples: []string{"/대화 내용:안녕", "/대화 내용:냠냠"},
	},
	Category:                   Chatting,
	RegisterApplicationCommand: true,
	RegisterMessageCommand:     false,
	ChatInputRun: func(ctx *ChatInputContext) {
		i := ctx.Inter
		i.DeferReply(&discordgo.InteractionResponseData{})

		result := chatbot.ParseResult(chatbot.ChatBot.GetResponse(i.Member.User.ID, i.Options["내용"].StringValue()), ctx.Inter.Session, i)
		i.EditReply(&utils.InteractionEdit{
			Content: &result,
		})
	},
}
