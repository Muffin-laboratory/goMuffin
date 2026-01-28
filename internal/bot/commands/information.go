package commands

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	subcommands "github.com/Muffin-laboratory/goMuffin/internal/bot/commands/subcommands/information"
	"github.com/bwmarrin/discordgo"
)

const (
	informationCommandBot        = "봇"
	informationCommandUser       = "유저"
	informationCommandPatchNotes = "패치내역"
)

var InformationCommand *Command = &Command{
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
	Flags: CommandFlagsIsBlocked,
	Run: func(inter *builders.InteractionCreate) error {
		switch inter.ApplicationCommandData().Options[0].Name {
		case informationCommandBot:
			if err := inter.DeferReply(nil); err != nil {
				return err
			}

			return subcommands.InfoBot(inter)
		case informationCommandUser:
			if err := inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			}); err != nil {
				return err
			}

			return subcommands.InfoUser(inter)
		case informationCommandPatchNotes:
			if err := inter.DeferReply(nil); err != nil {
				return err
			}

			return subcommands.InfoPatchLogs(inter)
		default:
			return nil
		}
	},
}

func init() {
	GetDiscommand().LoadCommand(InformationCommand)
}
