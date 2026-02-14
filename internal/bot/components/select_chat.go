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

var SelectChatComponent = &loader.Component{
	Middlewares: handler.Middlewares{
		middlewares.CheckIDMiddleware(),
		middlewares.TimeoutAndDeferMiddleware(loader.Timeout(), discord.InteractionTypeComponent, true, false),
	},
	Handle: func(r handler.Router) {
		r.Component(utils.SelectChat+"/{id}", func(e *handler.ComponentEvent) error {
			id, _ := bson.ObjectIDFromHex(e.Vars["id"])

			if _, err := repository.GetDatabase().Users.Update(e.Ctx, int64(e.User().ID), &repository.UserUpdate{
				ChatID: &id,
			}); err != nil {
				return err
			}

			_, err := e.UpdateInteractionResponse(
				discord.NewMessageUpdateV2([]discord.LayoutComponent{
					builders.MakeSuccessContainer("해당 채팅으로 변경하였어요."),
				}),
			)
			return err
		})
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(SelectChatComponent)
}
