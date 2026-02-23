package commands

import (
	"strings"

	subcommands "github.com/Muffin-laboratory/goMuffin/internal/bot/handler/commands/subcommands/information"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func init() {
	const (
		name                  = "정보"
		botCommandName        = "봇"
		userCommandName       = "유저"
		patchNotesCommandName = "패치내역"
	)

	options := []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        botCommandName,
			Description: "해당 봇의 정보를 확인해요.",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        userCommandName,
			Description: "명령어를 친 유저의 정보를 확인해요.",
		},
	}

	isOwnerEmpty := configs.Configs().GitHub.Owner == ""
	isRepoEmpty := configs.Configs().GitHub.Repository == ""

	if !isOwnerEmpty && !isRepoEmpty {
		options = append(options, discord.ApplicationCommandOptionSubCommand{
			Name:        patchNotesCommandName,
			Description: "해당 봇의 패치 내역을 확인해요.",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:         "버전",
					Description:  "패치 내역을 확인할 버전을 정해요.",
					Autocomplete: true,
				},
			},
		})
	}

	loader.GetDiscommand().RegisterCommand(discord.SlashCommandCreate{
		Name:        name,
		Description: "해당 봇의 정보를 알려줘요.",
		Options:     options,
	})

	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(
			middlewares.CheckBlocked(),
			middlewares.TimeoutAndDefer(loader.Timeout(), discord.InteractionTypeApplicationCommand, false, true),
		)

		r.Route("/"+name, func(r handler.Router) {
			r.Command("/"+botCommandName, subcommands.InfoBot)
			r.Command("/"+userCommandName, subcommands.InfoUser)

			if !isOwnerEmpty && !isRepoEmpty {
				r.Command("/"+patchNotesCommandName, subcommands.InfoPatchLogs)
				r.Autocomplete("/"+patchNotesCommandName, func(e *handler.AutocompleteEvent) error {
					var selectedTags []discord.AutocompleteChoice

					tags := repository.Tags()
					focused := e.Data.Focused().String()

					for i, tag := range tags {
						if focused == "" {
							selectedTags = append(selectedTags, discord.AutocompleteChoiceInt{
								Name:  tag,
								Value: i + 1,
							})
						} else {
							if strings.Contains(tag, focused) {
								selectedTags = append(selectedTags, discord.AutocompleteChoiceInt{
									Name:  tag,
									Value: i + 1,
								})
							}
						}
					}

					if len(selectedTags) > 25 {
						selectedTags = selectedTags[:25]
					}

					return e.AutocompleteResult(selectedTags)
				})
			}
		})
	})
}
