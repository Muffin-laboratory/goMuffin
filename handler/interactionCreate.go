package handler

import (
	"git.wh64.net/muffin/goMuffin/commands"
	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		commands.Discommand.ChatInputRun(i.ApplicationCommandData().Name, s, i)
		return
	case discordgo.InteractionMessageComponent:
		commands.Discommand.ComponentRun(s, i)
	case discordgo.InteractionModalSubmit:
		commands.Discommand.ModalRun(s, i)
	}
}
