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
		middlewares.TimeoutAndDeferMiddleware(loader.Timeout(), discord.InteractionTypeComponent, true, false),
	},
	Handle: func(r handler.Router) {
		r.Component(utils.ServiceAgree+"/{user_id}", func(e *handler.ComponentEvent) error {
			if e.User().ID.String() != e.Vars["user_id"] {
				_, err := e.CreateFollowupMessage(
					discord.NewMessageCreateBuilder().
						SetComponents(builders.MakeHasNoPermissionContainer()).
						SetIsComponentsV2(true).
						SetEphemeral(true).
						Build(),
				)
				return err
			}

			if _, err := repository.GetDatabase().Users.Create(e.Ctx, int64(e.User().ID)); err != nil {
				return err
			}

			bot, _ := e.Client().Caches.SelfUser()
			_, err := e.UpdateInteractionResponse(
				discord.NewMessageUpdateBuilder().
					SetComponents(builders.MakeSuccessContainer("가입을 했어요. 이제 %s의 모든 기능을 사용할 수 있어요.", bot.Username)).
					SetIsComponentsV2(true).
					Build(),
			)
			return err
		})

		r.Component(utils.ServiceDisagree+"/{user_id}", func(e *handler.ComponentEvent) error {
			if e.User().ID.String() != e.Vars["user_id"] {
				_, err := e.CreateFollowupMessage(
					discord.NewMessageCreateBuilder().
						SetComponents(builders.MakeHasNoPermissionContainer()).
						SetIsComponentsV2(true).
						SetEphemeral(true).
						Build(),
				)
				return err
			}

			_, err := e.UpdateInteractionResponse(
				discord.NewMessageUpdateBuilder().
					SetComponents(builders.MakeDeclineContainer("가입을 거부했어요.")).
					SetIsComponentsV2(true).
					Build(),
			)
			return err
		})
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(RegisterComponent)
}
