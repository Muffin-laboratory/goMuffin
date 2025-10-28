package commands

import (
	"context"

	"git.wh64.net/muffin/goMuffin/builders"
	subcommands "git.wh64.net/muffin/goMuffin/commands/subcommands/chat"
	"git.wh64.net/muffin/goMuffin/repository"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	chatCommandChatting = "하기"
	chatCommandList     = "목록"
	chatCommandCreate   = "생성"
	chatCommandDelete   = "삭제"
	chatCommandSettings = "설정"
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
	Flags: CommandFlagsIsRegistered | CommandFlagsIsBlocked,
	Run: func(inter *builders.InteractionCreate) error {
		switch opt := inter.ApplicationCommandData().Options[0]; opt.Name {
		case chatCommandChatting:
			if err := inter.DeferReply(nil); err != nil {
				return err
			}

			return subcommands.Chat(inter, builders.MakeCommandInteractionOptionsMap(opt.Options))
		case chatCommandCreate:
			if err := inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			}); err != nil {
				return err
			}

			return subcommands.Create(inter, builders.MakeCommandInteractionOptionsMap(opt.Options))
		case chatCommandList:
			if err := inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			}); err != nil {
				return err
			}

			return subcommands.List(inter)
		case chatCommandDelete:
			if err := inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			}); err != nil {
				return err
			}

			return subcommands.Delete(inter, builders.MakeCommandInteractionOptionsMap(opt.Options))
		case chatCommandSettings:
			return subcommands.Settings(inter)
		default:
			return nil
		}
	},
	Autocomplete: func(inter *builders.InteractionCreate) error {
		var choices []*discordgo.ApplicationCommandOptionChoice
		var data []*repository.Chat
		var focusedValue string

		for _, opt := range inter.ApplicationCommandData().Options[0].Options {
			if opt.Focused {
				focusedValue = opt.StringValue()
				break
			}
		}

		cur, err := repository.GetDatabase().Chats.Find(context.TODO(), bson.M{"name": bson.M{"$regex": focusedValue}})
		if err != nil {
			return err
		}

		defer cur.Close(context.TODO())

		if err = cur.All(context.TODO(), &data); err != nil {
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
	GetDiscommand().LoadCommand(ChatCommand)
}
