package chat

import (
	"fmt"
	"math/rand"

	"github.com/LoperLee/golang-hangul-toolkit/hangul"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func Create(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	name := fmt.Sprintf("새로운 채팅 %06d", rand.Intn(999999))

	if value, ok := data.OptString("이름"); ok {
		name = value
	}

	dbUser, err := repository.GetDatabase().Users.FindByID(e.Ctx, int64(e.User().ID))
	if err != nil {
		return err
	}

	if dbUser.ChattingMode == repository.ChattingMuffinMode {
		return chatSendErrorMessage(e)
	}

	if _, err := repository.GetDatabase().Chats.Create(e.Ctx, int64(e.User().ID), name); err != nil {
		return err
	}

	_, err = e.UpdateInteractionResponse(
		discord.NewMessageUpdateV2([]discord.LayoutComponent{
			builders.MakeSuccessContainer("%s 생성했어요. 이제 현재 채팅은 %s에요.", hangul.GetJosa(name, hangul.EUL_REUL), name),
		}),
	)
	return err
}
