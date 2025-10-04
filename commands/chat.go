package commands

import (
	subcommands "git.wh64.net/muffin/goMuffin/commands/subcommands/chat"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var (
	chatCommandChatting                  = "하기"
	chatCommandList                      = "목록"
	chatCommandCreate                    = "생성"
	chatCommandDelete                    = "삭제"
	chatCommandSwitchMode                = "모드전환"
	chatCommandReplyUser                 = "답장"
	chatCommandCreateNewChatAfter12Hours = "12시간"
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
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "모드",
						Description: "무슨 모드로 바꿀지 선택하세요.",
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "AI 모드",
								Value: databases.ChattingAIMode,
							},
							{
								Name:  "일반 모드",
								Value: databases.ChattingMuffinMode,
							},
						},
						Required: true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        chatCommandReplyUser,
				Description: "이 봇이 대답할 때 멘션을 킬지 선택해요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "활성화",
						Description: "활성화 여부를 선택해요.",
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "활성화",
								Value: 1,
							},
							{
								Name:  "비활성화",
								Value: 0,
							},
						},
						Required: true,
					},
				},
			},
		},
	},
	Category: Chatting,
	Flags:    CommandFlagsIsRegistered | CommandFlagsIsBlocked,
	Run: func(ctx *ChatInputContext) error {
		switch opt := ctx.Inter.ApplicationCommandData().Options[0]; opt.Name {
		case chatCommandChatting:
			if err := ctx.Inter.DeferReply(nil); err != nil {
				return err
			}

			return subcommands.Chat(ctx.Inter, utils.MakeCommandInteractionOptionsMap(opt.Options))
		case chatCommandCreate:
			if err := ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			}); err != nil {
				return err
			}

			return subcommands.Create(ctx.Inter, utils.MakeCommandInteractionOptionsMap(opt.Options))
		case chatCommandList:
			if err := ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			}); err != nil {
				return err
			}

			return subcommands.List(ctx.Inter)
		case chatCommandDelete:
			if err := ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			}); err != nil {
				return err
			}

			return subcommands.Delete(ctx.Inter, utils.MakeCommandInteractionOptionsMap(opt.Options))
		case chatCommandSwitchMode:
			if err := ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			}); err != nil {
				return err
			}

			return subcommands.SwitchMode(ctx.Inter, utils.MakeCommandInteractionOptionsMap(opt.Options))
		case chatCommandReplyUser:
			if err := ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			}); err != nil {
				return err
			}

			return subcommands.SetReplyUser(ctx.Inter, utils.MakeCommandInteractionOptionsMap(opt.Options))
		default:
			return nil
		}
	},
}

func init() {
	GetDiscommand().LoadCommand(ChatCommand)
}
