package handler

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var err error
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		err = commands.Discommand.ChatInputRun(i.ApplicationCommandData().Name, s, i)
		if err != nil {
			goto ErrMsg
		}
	case discordgo.InteractionMessageComponent:
		err = commands.Discommand.ComponentRun(s, i)
		if err != nil {
			goto ErrMsg
		}
	case discordgo.InteractionModalSubmit:
		err = commands.Discommand.ModalRun(s, i)
		if err != nil {
			goto ErrMsg
		}
	}

	// 아 몰라 goto 쓸래
ErrMsg:
	owner, _ := s.User(configs.Config.Bot.OwnerId)
	utils.NewMessageSender(&utils.InteractionCreate{
		InteractionCreate: i,
		Session:           s,
	}).
		AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: fmt.Sprintf("오류가 발생하였어요. 만약 계속 발생한다면, %s으로 연락해주세요.", utils.InlineCode(owner.Username))})).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
