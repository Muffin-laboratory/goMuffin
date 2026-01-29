package commands

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	subcommands "github.com/Muffin-laboratory/goMuffin/internal/bot/commands/subcommands/information"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/bwmarrin/discordgo"
)

const (
	informationCommandBot        = "봇"
	informationCommandUser       = "유저"
	informationCommandPatchNotes = "패치내역"
)

var InformationCommand = &loader.Command{
	Deferred: true,
	DeferOptions: &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
	},
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "정보",
		Description: "해당 봇의 정보를 알려줘요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        informationCommandBot,
				Description: "해당 봇의 정보를 확인해요.",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        informationCommandUser,
				Description: "명령어를 친 유저의 정보를 확인해요.",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        informationCommandPatchNotes,
				Description: "해당 봇의 패치 내역을 확인해요.",
			},
		},
	},
	Flags: loader.CommandFlagsIsBlocked,
	Run: func(inter *builders.InteractionCreate) error {
		switch inter.ApplicationCommandData().Options[0].Name {
		case informationCommandBot:
			return subcommands.InfoBot(inter)
		case informationCommandUser:
			return subcommands.InfoUser(inter)
		case informationCommandPatchNotes:
			return subcommands.InfoPatchLogs(inter)
		default:
			return nil
		}
	},
}

func init() {
	loader.GetDiscommand().LoadCommand(InformationCommand)
}
