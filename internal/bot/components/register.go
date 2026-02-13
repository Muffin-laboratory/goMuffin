package components

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var RegisterComponent = &loader.Component{
	Middlewares: handler.Middlewares{
		middlewares.CheckIDMiddleware(),
		middlewares.TimeoutAndDeferMiddleware(loader.Timeout(), discord.InteractionTypeComponent, true, false),
	},
	Handle: func(r handler.Router) {
		r.Component(utils.ServiceAgree, func(e *handler.ComponentEvent) error {
			if _, err := repository.GetDatabase().Users.Create(e.Ctx, int64(e.User().ID)); err != nil {
				return err
			}

			bot, _ := e.Client().Caches.SelfUser()
			_, err := e.UpdateInteractionResponse(
				discord.NewMessageUpdateV2([]discord.LayoutComponent{
					builders.MakeSuccessContainer("가입을 했어요. 이제 %s의 모든 기능을 사용할 수 있어요.", bot.Username),
				}),
			)
			return err
		})

		r.Component(utils.ServiceDisagree, func(e *handler.ComponentEvent) error {
			_, err := e.UpdateInteractionResponse(
				discord.NewMessageUpdateV2([]discord.LayoutComponent{
					builders.MakeDeclineContainer("가입을 거부했어요."),
				}),
			)
			return err
		})
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(RegisterComponent)
}
