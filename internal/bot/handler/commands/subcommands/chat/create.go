package chat

import (
	"fmt"
	"math/rand"

	"github.com/LoperLee/golang-hangul-toolkit/hangul"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders/customid"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func Create(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	dbUser, err := repository.GetDatabase().Users.FindByID(e.Ctx, int64(e.User().ID))
	if err != nil {
		return err
	}

	if dbUser.ChattingMode == repository.ChattingMuffinMode {
		return chatSendErrorMessage(e)
	}

	name := fmt.Sprintf("새로운 채팅 %06d", rand.Intn(999999))

	if value, ok := data.OptString("이름"); ok {
		name = value
	}

	if data.Bool("프롬프트_지정") {
		return e.Modal(
			discord.NewModalCreate(
				customid.MakeCreateChat(name),
				"사용자 지정 프롬프트 설정",
				[]discord.LayoutComponent{
					discord.NewLabel(
						"사용자 지정 프롬프트",
						discord.NewParagraphTextInput(customid.CreateChatSetPrompt).
							WithPlaceholder("여기에 프롬프트를 입력..."),
					).
						WithDescription(
							"여기에 사용자 지정 프롬프트를 지정해주세요. " +
								"단, 공란이면 유저의 설정한 프롬프트로 설정돼요.",
						),
				},
			),
		)
	}

	err = e.DeferCreateMessage(true)
	if err != nil {
		return err
	}

	if _, err := repository.GetDatabase().Chats.Create(e.Ctx, int64(e.User().ID), dbUser.Prompt, name); err != nil {
		return err
	}

	_, err = e.UpdateInteractionResponse(
		discord.NewMessageUpdateV2([]discord.LayoutComponent{
			builders.MakeSuccessContainer("%s 생성했어요. 이제 현재 채팅은 %s에요.", hangul.GetJosa(name, hangul.EUL_REUL), name),
		}),
	)
	return err
}
