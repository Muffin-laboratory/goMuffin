package components

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var DeleteKnowledgeComponent = &loader.Component{
	Middlewares: handler.Middlewares{
		middlewares.CheckIDMiddleware(),
		middlewares.TimeoutAndDeferMiddleware(loader.Timeout(), discord.InteractionTypeComponent, true, false),
	},
	Handle: func(r handler.Router) {
		r.Component(utils.DeleteKnowledge+"/{id}", func(e *handler.ComponentEvent) error {
			data := e.Vars["id"]

			id, _ := bson.ObjectIDFromHex(data)

			if err := repository.GetDatabase().Knowledge.DeleteByID(e.Ctx, id); err != nil {
				return err
			}

			_, err := e.UpdateInteractionResponse(
				discord.NewMessageUpdateV2([]discord.LayoutComponent{
					builders.MakeSuccessContainer("해당 항목을 삭제했어요."),
				}),
			)
			return err
		})
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(DeleteKnowledgeComponent)
}
