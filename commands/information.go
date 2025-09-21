package commands

import (
	subcommands "git.wh64.net/muffin/goMuffin/commands/subcommands/information"
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
	Category: General,
	Flags:    CommandFlagsIsBlocked,
	Run: func(ctx *ChatInputContext) error {
		switch ctx.Inter.ApplicationCommandData().Options[0].Name {
		case informationCommandBot:
			if err := ctx.Inter.DeferReply(nil); err != nil {
				return err
			}

			return subcommands.InfoBot(ctx.Inter)
		case informationCommandUser:
			if err := ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			}); err != nil {
				return err
			}

			return subcommands.InfoUser(ctx.Inter)
		case informationCommandPatchNotes:
			if err := ctx.Inter.DeferReply(nil); err != nil {
				return err
			}

			return subcommands.InfoPatchLogs(ctx.Inter)
		default:
			return nil
		}
	},
}

func init() {
	GetDiscommand().LoadCommand(InformationCommand)
}
