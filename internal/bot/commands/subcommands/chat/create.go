package chat

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
)

func Create(ctx context.Context, i *builders.CommandCreate) error {
	name := fmt.Sprintf("새로운 채팅 %06d", rand.Intn(999999))
	commandData := i.SlashCommandInteractionData()

	if value, ok := commandData.OptString("이름"); ok {
		name = value
	}

	dbUser, err := repository.GetDatabase().Users.FindByID(ctx, i.User().ID.String())
	if err != nil {
		return err
	}

	if dbUser.ChattingMode == repository.ChattingMuffinMode {
		return chatSendErrorMessage(i)
	}

	if _, err := repository.GetDatabase().Chats.Create(ctx, i.User().ID.String(), name); err != nil {
		return err
	}

	return builders.NewMessageSender(i).
		AddComponents(builders.MakeSuccessContainer("%s를 생성했어요. 이제 현재 채팅은 %s에요.", name, name)).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
