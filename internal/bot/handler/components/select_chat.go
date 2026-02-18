package components

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders/customid"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func init() {
	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(
			middlewares.CheckIDMiddleware(),
			middlewares.TimeoutAndDeferMiddleware(loader.Timeout(), discord.InteractionTypeComponent, true, false),
		)

		r.Component(customid.SelectChat+"/{id}", func(e *handler.ComponentEvent) error {
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
	})
}
