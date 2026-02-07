package commands

import (
	"context"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	subcommands "github.com/Muffin-laboratory/goMuffin/internal/bot/commands/subcommands/information"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/disgoorg/disgo/discord"
)

const (
	informationCommandBot        = "봇"
	informationCommandUser       = "유저"
	informationCommandPatchNotes = "패치내역"
)

var InformationCommand = &loader.Command{
	Deferred:         true,
	IsDeferEphemeral: true,
	SlashCommandCreate: &discord.SlashCommandCreate{
		Name:        "정보",
		Description: "해당 봇의 정보를 알려줘요.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        informationCommandBot,
				Description: "해당 봇의 정보를 확인해요.",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        informationCommandUser,
				Description: "명령어를 친 유저의 정보를 확인해요.",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        informationCommandPatchNotes,
				Description: "해당 봇의 패치 내역을 확인해요.",
			},
		},
	},
	Flags: loader.CommandFlagsIsBlocked,
	Run: func(ctx context.Context, inter *builders.CommandCreate) error {
		switch *inter.SlashCommandInteractionData().SubCommandName {
		case informationCommandBot:
			return subcommands.InfoBot(ctx, inter)
		case informationCommandUser:
			return subcommands.InfoUser(ctx, inter)
		case informationCommandPatchNotes:
			return subcommands.InfoPatchLogs(ctx, inter)
		default:
			return nil
		}
	},
}

func init() {
	loader.GetDiscommand().LoadCommand(InformationCommand)
}
