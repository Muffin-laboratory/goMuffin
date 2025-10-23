package chat

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/repository"
)

func SwitchMode(i *builders.InteractionCreate, opts builders.CommandInteractionOptionsMap) error {
	newMode := repository.ChattingMode(opts["모드"].IntValue())

	if _, err := repository.GetDatabase().Users.Update(i.User.ID, &repository.UserUpdate{
		ChattingMode: &newMode,
	}); err != nil {
		return err
	}

	return builders.NewMessageSender(i).
		AddComponents(builders.MakeSuccessContainer(fmt.Sprintf("모드를 성공적으로 %s로 바꿨어요.", repository.ModeString(newMode)))).
		SetComponentsV2(true).
		Send()
}
