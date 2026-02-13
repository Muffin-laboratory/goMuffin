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
	"go.mongodb.org/mongo-driver/v2/bson"
)

var DeleteKnowledgeComponent = &loader.Component{
	Middlewares: handler.Middlewares{
		middleware.Defer(discord.InteractionTypeComponent, true, false),
		middlewares.TimeoutMiddleware(loader.Timeout()),
	},
	Handle: func(r handler.Router) {
		r.Component(utils.DeleteKnowledge+"/{id}/{user_id}", func(e *handler.ComponentEvent) error {
			data := e.Vars["id"]

			if e.User().ID.String() != e.Vars["user_id"] {
				_, err := e.UpdateInteractionResponse(
					discord.NewMessageUpdateBuilder().
						SetComponents(builders.MakeHasNoPermissionContainer()).
						SetIsComponentsV2(true).
						Build(),
				)
				return err
			}

			id, _ := bson.ObjectIDFromHex(data)

			if err := repository.GetDatabase().Knowledge.DeleteByID(e.Ctx, id); err != nil {
				return err
			}

			_, err := e.UpdateInteractionResponse(
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
