package commands

import (
	"context"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	subcommands "github.com/Muffin-laboratory/goMuffin/internal/bot/commands/subcommands/chat"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var (
	chatCommandChatting = "하기"
	chatCommandList     = "목록"
	chatCommandCreate   = "생성"
	chatCommandDelete   = "삭제"
	chatCommandSettings = "설정"
)

var chatNameMaxLength = 25

var ChatCommand = &loader.Command{
	Deferred:         true,
	IsDeferEphemeral: true,
	SlashCommandCreate: &discord.SlashCommandCreate{
		Name:        "대화",
		Description: "이 봇이랑 대화해요.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        chatCommandList,
				Description: "채팅 목록을 나열해요.",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        chatCommandCreate,
				Description: "새로운 채팅을 생성해요.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "이름",
						Description: "채팅방의 이름 (25자 이내)",
						MaxLength:   &chatNameMaxLength,
						Required:    false,
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        chatCommandChatting,
				Description: "이 봇이랑 대화해요.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "내용",
						Description: "대화할 내용",
						Required:    true,
					},
					discord.ApplicationCommandOptionAttachment{
						Name:        "첨부파일",
						Description: "같이 보낼 파일",
						Required:    false,
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        chatCommandDelete,
				Description: "채팅을 삭제해요.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:         "이름",
						Description:  "지울 채팅방의 이름",
						MaxLength:    &chatNameMaxLength,
						Required:     true,
						Autocomplete: true,
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        chatCommandSettings,
				Description: "봇의 설정을 개인화 해요.",
			},
		},
	},
	Flags: loader.CommandFlagsIsRegistered | loader.CommandFlagsIsBlocked,
	Run: func(ctx context.Context, inter *builders.CommandCreate) error {

		switch *inter.SlashCommandInteractionData().SubCommandName {
		case chatCommandChatting:
			return subcommands.Chat(ctx, inter)
		case chatCommandCreate:
			return subcommands.Create(ctx, inter)
		case chatCommandList:
			return subcommands.List(ctx, inter)
		case chatCommandDelete:
			return subcommands.Delete(ctx, inter)
		case chatCommandSettings:
			return subcommands.Settings(ctx, inter)
		default:
			return nil
		}
	},
	Autocomplete: func(ctx context.Context, inter *events.AutocompleteInteractionCreate) error {
		var choices []discord.AutocompleteChoice

		focusedValue := inter.Data.Focused().String()

		filter := query.ChatQueryBuilder().SetNameByRegex(focusedValue)
		data, err := repository.GetDatabase().Chats.Find(ctx, filter)
		if err != nil {
			return err
		}

		if len(data) > 25 {
			data = data[:25]
		}

		for _, data := range data {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  data.Name,
				Value: data.Name,
			})
		}

		return inter.AutocompleteResult(choices)
	},
}

func init() {
	loader.GetDiscommand().LoadCommand(ChatCommand)
}
