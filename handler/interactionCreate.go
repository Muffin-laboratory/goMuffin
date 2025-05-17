package handler

import (
	"git.wh64.net/muffin/goMuffin/commands"
	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		commands.Discommand.ChatInputRun(i.ApplicationCommandData().Name, s, i)
		return
	} else if i.Type == discordgo.InteractionMessageComponent {
		commands.Discommand.ComponentRun(s, i)
	} else if i.Type == discordgo.InteractionModalSubmit {
		commands.Discommand.ModalRun(s, i)
	}
}
