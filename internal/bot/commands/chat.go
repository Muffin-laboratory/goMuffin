package commands

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	subcommands "github.com/Muffin-laboratory/goMuffin/internal/bot/commands/subcommands/chat"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/bwmarrin/discordgo"
)

var (
	chatCommandChatting = "하기"
	chatCommandList     = "목록"
	chatCommandCreate   = "생성"
	chatCommandDelete   = "삭제"
	chatCommandSettings = "설정"
)

const chatNameMaxLength = 25

var ChatCommand = &loader.Command{
	Deferred: true,
	DeferOptions: &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
	},
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
						Required:    false,
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
					{
						Type:        discordgo.ApplicationCommandOptionAttachment,
						Name:        "첨부파일",
						Description: "같이 보낼 파일",
						Required:    false,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        chatCommandDelete,
				Description: "채팅을 삭제해요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "이름",
						Description:  "지울 채팅방의 이름",
						MaxLength:    chatNameMaxLength,
						Required:     true,
						Autocomplete: true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        chatCommandSettings,
				Description: "봇의 설정을 개인화 해요.",
			},
		},
	},
	Flags: loader.CommandFlagsIsRegistered | loader.CommandFlagsIsBlocked,
	Run: func(inter *builders.InteractionCreate) error {
		switch opt := inter.ApplicationCommandData().Options[0]; opt.Name {
		case chatCommandChatting:
			return subcommands.Chat(inter, builders.MakeCommandInteractionOptionsMap(opt.Options))
		case chatCommandCreate:
			return subcommands.Create(inter, builders.MakeCommandInteractionOptionsMap(opt.Options))
		case chatCommandList:
			return subcommands.List(inter)
		case chatCommandDelete:
			return subcommands.Delete(inter, builders.MakeCommandInteractionOptionsMap(opt.Options))
		case chatCommandSettings:
			return subcommands.Settings(inter)
		default:
			return nil
		}
	},
	Autocomplete: func(inter *builders.InteractionCreate) error {
		var choices []*discordgo.ApplicationCommandOptionChoice
		var focusedValue string

		for _, opt := range inter.ApplicationCommandData().Options[0].Options {
			if opt.Focused {
				focusedValue = opt.StringValue()
				break
			}
		}

		data, err := repository.GetDatabase().Chats.Find(inter.Ctx, query.ChatQueryBuilder().SetNameByRegex(focusedValue))
		if err != nil {
			return err
		}

		for _, data := range data {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  data.Name,
				Value: data.Name,
			})
		}

		return inter.Autocomplete(choices)
	},
}

func init() {
	loader.GetDiscommand().LoadCommand(ChatCommand)
}
