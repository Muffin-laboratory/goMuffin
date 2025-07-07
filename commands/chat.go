package commands

import (
	"log"

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
	Flags:                      CommandFlagsIsRegistered | CommandFlagsIsBlocked,
	ChatInputRun: func(ctx *ChatInputContext) error {
		i := ctx.Inter
		i.DeferReply(&discordgo.InteractionResponseData{})

		str, err := chatbot.ChatBot.GetResponse(i.Member.User, i.Options["내용"].StringValue())
		if err != nil {
			log.Println(err)
			i.EditReply(&utils.InteractionEdit{
				Content: &str,
			})
			return nil
		}

		result := chatbot.ParseResult(str, ctx.Inter.Session, i)
		return i.EditReply(&utils.InteractionEdit{
			Content: &result,
		})
	},
}
