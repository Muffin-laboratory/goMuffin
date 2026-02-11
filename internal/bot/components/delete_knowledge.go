package components

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"
)

var DeleteKnowledgeComponent = &loader.Component{
	Middlewares: handler.Middlewares{
		middleware.Defer(discord.InteractionTypeComponent, true, false),
		middlewares.TimeoutMiddleware(loader.Timeout()),
	},
	Handle: func(r handler.Router) {
		r.Component(utils.DeleteKnowledge+"/{data}", func(inter *handler.ComponentEvent) error {
			data := inter.Vars["data"]

			userID := utils.GetDeleteKnowledgeUserID(data)
			if inter.User().ID.String() != userID {
				_, err := inter.UpdateInteractionResponse(
					discord.NewMessageUpdateBuilder().
						SetComponents(builders.MakeHasNoPermissionContainer()).
						SetIsComponentsV2(true).
						Build(),
				)
				return err
			}

			id := utils.GetDeleteKnowledgeID(data)

			if err := repository.GetDatabase().Knowledge.DeleteByID(inter.Ctx, id); err != nil {
				return err
			}

			_, err := inter.UpdateInteractionResponse(
				discord.NewMessageUpdateBuilder().
					SetComponents(builders.MakeSuccessContainer("해당 항목을 삭제했어요.")).
					SetIsComponentsV2(true).
					Build(),
			)
			return err
		})
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(DeleteKnowledgeComponent)
}
