package commands

import (
	"context"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	subcommands "github.com/Muffin-laboratory/goMuffin/internal/bot/commands/subcommands/knowledge"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

const (
	knowledgeLearn  = "추가"
	knowledgeList   = "목록"
	knowledgeDelete = "삭제"
)

var knowledgeMaxCommandLength = 100

var KnowledgeCommand = &loader.Command{
	Deferred:         true,
	IsDeferEphemeral: true,
	SlashCommandCreate: &discord.SlashCommandCreate{
		Name:        "지식",
		Description: "이 봇이 사용자와 대화할 때 알면 좋은 지식을 관리하는 명령어에요.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        knowledgeLearn,
				Description: "단어를 가르치는 명령어에요.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "단어",
						Description: "등록할 단어",
						Required:    true,
						MaxLength:   &knowledgeMaxCommandLength,
					},
					discord.ApplicationCommandOptionString{
						Name:        "대답",
						Description: "해당 단어의 대답",
						Required:    true,
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        knowledgeList,
				Description: "당신이 가르쳐준 지식을 나열해요.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:         "단어",
						Description:  "해당 단어가 포함된 결과",
						Required:     false,
						Autocomplete: true,
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        chatCommandDelete,
				Description: "당신이 가르쳐준 단어를 삭제해요.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:         "단어",
						Description:  "삭제할 단어",
						Required:     true,
						Autocomplete: true,
					},
				},
			},
		},
	},
	Flags: loader.CommandFlagsIsRegistered | loader.CommandFlagsIsBlocked,
	Run: func(ctx context.Context, inter *builders.CommandCreate) error {
		switch *inter.SlashCommandInteractionData().SubCommandName {
		case knowledgeLearn:
			igCommands := []string{}

			for _, command := range loader.GetDiscommand().Commands {
				igCommands = append(igCommands, command.Name)
			}

			return subcommands.Learn(ctx, inter, igCommands)
		case knowledgeList:
			return subcommands.List(ctx, inter)
		case knowledgeDelete:
			return subcommands.Delete(ctx, inter)
		default:
			return nil
		}
	},
	Autocomplete: func(ctx context.Context, inter *events.AutocompleteInteractionCreate) error {
		var choices []discord.AutocompleteChoice
		var focusedValue string

		for _, opt := range inter.Data.Options {
			if opt.Focused {
				focusedValue = opt.String()
				break
			}
		}

		filter := query.KnowledgeQueryBuilder().SetUserID(int64(inter.User().ID)).SetCommandByRegex(focusedValue)
		data, err := repository.GetDatabase().Knowledge.Find(ctx, filter)
		if err != nil {
			return err
		}

		for _, data := range data {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  data.Command,
				Value: data.Command,
			})
		}

		return inter.AutocompleteResult(choices)
	},
}

func init() {
	loader.GetDiscommand().LoadCommand(KnowledgeCommand)
}
