package commands

import (
	subcommands "github.com/Muffin-laboratory/goMuffin/internal/bot/handler/commands/subcommands/knowledge"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func init() {
	var commandMaxLength = 100

	const (
		name              = "지식"
		learnCommandName  = "추가"
		listCommandName   = "목록"
		deleteCommandName = "삭제"
	)

	loader.GetDiscommand().RegisterCommand(discord.SlashCommandCreate{
		Name:        name,
		Description: "이 봇이 사용자와 대화할 때 알면 좋은 지식을 관리하는 명령어에요.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        learnCommandName,
				Description: "단어를 가르치는 명령어에요.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "단어",
						Description: "등록할 단어",
						Required:    true,
						MaxLength:   &commandMaxLength,
					},
					discord.ApplicationCommandOptionString{
						Name:        "대답",
						Description: "해당 단어의 대답",
						Required:    true,
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        listCommandName,
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
				Name:        deleteCommandName,
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
	})

	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(
			middlewares.CheckUserAndBlockedMiddleware(),
			middlewares.TimeoutAndDeferMiddleware(loader.Timeout(), discord.InteractionTypeApplicationCommand, false, true),
		)

		r.Autocomplete("/"+name, func(e *handler.AutocompleteEvent) error {
			var choices []discord.AutocompleteChoice
			var focusedValue string

			for _, opt := range e.Data.Options {
				if opt.Focused {
					focusedValue = opt.String()
					break
				}
			}

			filter := query.KnowledgeQueryBuilder().SetUserID(int64(e.User().ID)).SetCommandByRegex(focusedValue)
			data, err := repository.GetDatabase().Knowledge.Find(e.Ctx, filter)
			if err != nil {
				return err
			}

			if len(data) > 25 {
				data = data[:25]
			}

			for _, data := range data {
				choices = append(choices, discord.AutocompleteChoiceString{
					Name:  data.Command,
					Value: data.Command,
				})
			}

			return e.AutocompleteResult(choices)
		})

		r.Route("/"+name, func(r handler.Router) {
			r.SlashCommand("/"+learnCommandName, subcommands.Learn)
			r.SlashCommand("/"+listCommandName, subcommands.List)
			r.SlashCommand("/"+deleteCommandName, subcommands.Delete)
		})
	})
}
