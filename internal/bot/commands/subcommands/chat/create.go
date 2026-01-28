package chat

import (
	"fmt"
	"math/rand"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
)

func Create(i *builders.InteractionCreate, opts builders.CommandInteractionOptionsMap) error {
	var name = fmt.Sprintf("새로운 채팅 %06d", rand.Intn(999999))

	if opt, ok := opts["이름"]; ok {
		name = opt.StringValue()
	}

	dbUser, err := repository.GetDatabase().Users.Get(i.User.ID)
	if err != nil {
		return err
	}

	if dbUser.ChattingMode == repository.ChattingMuffinMode {
		return chatSendErrorMessage(i)
	}

	if _, err := repository.GetDatabase().Chats.Create(i.User.ID, name); err != nil {
		return err
	}

	return builders.NewMessageSender(i).
		AddComponents(builders.MakeSuccessContainer(fmt.Sprintf("%s를 생성했어요. 이제 현재 채팅은 %s에요.", name, name))).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
