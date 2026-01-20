package handler

import (
	"fmt"
	"log"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if err := commands.GetDiscommand().ChatInputRun(i.ApplicationCommandData().Name, s, i); err != nil {
			returnErr(s, i, err)
		}
	case discordgo.InteractionMessageComponent:
		if err := commands.GetDiscommand().ComponentRun(s, i); err != nil {
			returnErr(s, i, err)
		}
	case discordgo.InteractionModalSubmit:
		if err := commands.GetDiscommand().ModalRun(s, i); err != nil {
			returnErr(s, i, err)
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		if err := commands.GetDiscommand().ChatInputAutocomplete(i.ApplicationCommandData().Name, s, i); err != nil {
			returnErr(s, i, err)
		}
	}
}

func returnErr(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	owner, _ := s.User(configs.GetConfig().Bot.OwnerID)
	builders.NewMessageSender(&builders.InteractionCreate{
		InteractionCreate: i,
		Session:           s,
	}).
		AddComponents(builders.MakeErrorContainer(fmt.Sprintf("오류가 발생하였어요. 만약 계속 발생한다면, %s으로 연락해주세요.", utils.InlineCode(owner.Username)))).
		SetComponentsV2(true).
		SetReply(true).
		Send()
	log.Println(err)
}
