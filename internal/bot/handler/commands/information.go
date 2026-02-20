package commands

import (
	subcommands "github.com/Muffin-laboratory/goMuffin/internal/bot/handler/commands/subcommands/information"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
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

	loader.GetDiscommand().RegisterCommand(discord.SlashCommandCreate{
		Name:        name,
		Description: "해당 봇의 정보를 알려줘요.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        botCommandName,
				Description: "해당 봇의 정보를 확인해요.",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        userCommandName,
				Description: "명령어를 친 유저의 정보를 확인해요.",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        patchNotesCommandName,
				Description: "해당 봇의 패치 내역을 확인해요.",
			},
		},
	})

	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(
			middlewares.CheckBlocked(),
			middlewares.TimeoutAndDefer(loader.Timeout(), discord.InteractionTypeApplicationCommand, false, true),
		)

		r.Route("/"+name, func(r handler.Router) {
			r.Command("/"+botCommandName, subcommands.InfoBot)
			r.Command("/"+userCommandName, subcommands.InfoUser)
			r.Command("/"+patchNotesCommandName, subcommands.InfoPatchLogs)
		})
	})
}
