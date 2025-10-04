package chat

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func SwitchMode(i *utils.InteractionCreate, opts utils.CommandInteractionOptionsMap) error {
	newMode := databases.ChattingMode(opts["모드"].IntValue())

	if _, err := databases.GetDatabase().Users.Update(i.User.ID, &databases.UserUpdate{
		ChattingMode: &newMode,
	}); err != nil {
		return err
	}

	return utils.NewMessageSender(i).
		AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("모드를 성공적으로 %s로 바꿨어요.", databases.ModeString(newMode))})).
		SetComponentsV2(true).
		Send()
}
