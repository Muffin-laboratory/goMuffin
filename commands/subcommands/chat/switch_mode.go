package chat

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/databases"
)

func SwitchMode(i *builders.InteractionCreate, opts builders.CommandInteractionOptionsMap) error {
	newMode := databases.ChattingMode(opts["모드"].IntValue())

	if _, err := databases.GetDatabase().Users.Update(i.User.ID, &databases.UserUpdate{
		ChattingMode: &newMode,
	}); err != nil {
		return err
	}

	return builders.NewMessageSender(i).
		AddComponents(builders.MakeSuccessContainer(fmt.Sprintf("모드를 성공적으로 %s로 바꿨어요.", databases.ModeString(newMode)))).
		SetComponentsV2(true).
		Send()
}
