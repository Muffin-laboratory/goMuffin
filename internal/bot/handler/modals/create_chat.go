package modals

import (
	"github.com/LoperLee/golang-hangul-toolkit/hangul"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders/customid"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/chatbot"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func init() {
	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(middlewares.TimeoutAndDefer(loader.Timeout(), discord.InteractionTypeModalSubmit, false, true))

		r.Modal(customid.CreateChat+"/{name}", func(e *handler.ModalEvent) error {
			name := e.Vars["name"]

			var prompt string

			dbUser, err := repository.GetDatabase().Users.FindByID(e.Ctx, int64(e.User().ID))
			if err != nil {
				return err
			}

			if value, ok := e.Data.OptText(customid.CreateChatSetPrompt); ok {
				prompt = value
			} else {
				if dbUser.Prompt != "" {
					prompt = dbUser.Prompt
				} else {
					prompt = chatbot.GetChatBot().GetPrompt()
				}
			}

			if _, err := repository.GetDatabase().Chats.Create(e.Ctx, int64(e.User().ID), prompt, name); err != nil {
				return err
			}

			_, err = e.UpdateInteractionResponse(
				discord.NewMessageUpdateV2([]discord.LayoutComponent{
					builders.MakeSuccessContainer("%s 생성했어요. 이제 현재 채팅은 %s에요.", hangul.GetJosa(name, hangul.EUL_REUL), name),
				}),
			)
			return err
		})
	})
}
