package commands

import (
	subcommand "git.wh64.net/muffin/goMuffin/commands/subcommands/chat"
	"github.com/bwmarrin/discordgo"
)

type chatCommandType string

var (
	chatCommandChatting   chatCommandType = "하기"
	chatCommandList       chatCommandType = "목록"
	chatCommandCreate     chatCommandType = "생성"
	chatCommandDelete     chatCommandType = "삭제"
	chatCommandSwitchMode chatCommandType = "모드전환"
)

const chatNameMaxLength = 25

var ChatCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "대화",
		Description: "이 봇이랑 대화해요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        string(chatCommandList),
				Description: "채팅 목록을 나열해요.",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        string(chatCommandCreate),
				Description: "새로운 채팅을 생성해요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "이름",
						Description: "채팅방의 이름 (25자 이내)",
						MaxLength:   chatNameMaxLength,
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        string(chatCommandChatting),
				Description: "이 봇이랑 대화해요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "내용",
						Description: "대화할 내용",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        string(chatCommandDelete),
				Description: "채팅을 삭제해요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "이름",
						Description: "지울 채팅방의 이름",
						MaxLength:   chatNameMaxLength,
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        string(chatCommandSwitchMode),
				Description: "채팅 방식을 변경해요. (일반 <-> AI)",
			},
		},
	},
	DetailedDescription: DetailedDescription{
		Usage: "/대화 (목록/생성/삭제) (이름:숫자(최대 25자, 목록에선 사용 불가능))",
		Examples: []string{
			"/대화 목록",
			"/대화 생성 이름:머핀 냠냠",
			"/대화 삭제 이름:뷁",
		},
	},
	Category: Chatting,
	Flags:    CommandFlagsIsRegistered | CommandFlagsIsBlocked,
	Run: func(ctx *ChatInputContext) error {
		if opt, ok := ctx.Inter.Options[string(chatCommandChatting)]; ok {
			ctx.Inter.DeferReply(nil)
			return subcommand.Chat(ctx.Inter, opt.Options[0].StringValue())
		} else if opt, ok := ctx.Inter.Options[string(chatCommandCreate)]; ok {
			ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			})
			return subcommand.Create(ctx.Inter, opt.Options[0].StringValue())
		} else if _, ok := ctx.Inter.Options[string(chatCommandList)]; ok {
			ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			})
			return subcommand.List(ctx.Inter)
		} else if _, ok := ctx.Inter.Options[string(chatCommandDelete)]; ok {
			ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			})
			return subcommand.Delete(ctx.Inter, opt.Options[0].StringValue())
		} else if _, ok := ctx.Inter.Options[string(chatCommandSwitchMode)]; ok {
			ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			})
			return subcommand.SwitchMode(ctx.Inter)
		}

		return nil
	},
}
