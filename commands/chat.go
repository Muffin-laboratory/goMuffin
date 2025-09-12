package commands

import (
	subcommands "git.wh64.net/muffin/goMuffin/commands/subcommands/chat"
	"github.com/bwmarrin/discordgo"
)

var (
	chatCommandChatting   = "하기"
	chatCommandList       = "목록"
	chatCommandCreate     = "생성"
	chatCommandDelete     = "삭제"
	chatCommandSwitchMode = "모드전환"
)

const chatNameMaxLength = 25

var ChatCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "대화",
		Description: "이 봇이랑 대화해요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        chatCommandList,
				Description: "채팅 목록을 나열해요.",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        chatCommandCreate,
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
				Name:        chatCommandChatting,
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
				Name:        chatCommandDelete,
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
				Name:        chatCommandSwitchMode,
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
		switch opt := ctx.Inter.ApplicationCommandData().Options[0]; opt.Name {
		case chatCommandChatting:
			ctx.Inter.DeferReply(nil)
			return subcommands.Chat(ctx.Inter, opt.Options[0].StringValue())
		case chatCommandSwitchMode:
			ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			})
			return subcommands.SwitchMode(ctx.Inter)
		case chatCommandCreate:
			ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			})
			return subcommands.Create(ctx.Inter, opt.Options[0].StringValue())
		case chatCommandList:
			ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			})
			return subcommands.List(ctx.Inter)
		case chatCommandDelete:
			ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			})
			return subcommands.Delete(ctx.Inter, opt.Options[0].StringValue())
		default:
			return nil
		}
	},
}

func init() {
	GetDiscommand().LoadCommand(ChatCommand)
}
