package handler

import (
	"log"

	"git.wh64.net/muffin/goMuffin/commands"
	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		err := commands.GetDiscommand().ChatInputRun(i.ApplicationCommandData().Name, s, i)
		if err != nil {
			log.Println(err)
		}
	case discordgo.InteractionMessageComponent:
		err := commands.GetDiscommand().ComponentRun(s, i)
		if err != nil {
			log.Println(err)
		}
	case discordgo.InteractionModalSubmit:
		err := commands.GetDiscommand().ModalRun(s, i)
		if err != nil {
			log.Println(err)
		}
	}
}
