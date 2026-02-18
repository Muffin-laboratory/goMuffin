package components

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders/customid"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
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

		r.Component(customid.DeleteChat+"/{value}", func(e *handler.ComponentEvent) error {
			value := e.Vars["value"]

			if value == "cancel" {
				_, err := e.UpdateInteractionResponse(
					discord.NewMessageUpdateV2([]discord.LayoutComponent{
						builders.MakeCanceledContainer("아무 채팅방을 삭제하지 않았어요."),
					}),
				)
				return err
			}

			id, _ := bson.ObjectIDFromHex(value)

			if err := repository.GetDatabase().Chats.DeleteByID(e.Ctx, id); err != nil {
				return err
			}

			if err := repository.GetDatabase().Memory.DeleteMany(e.Ctx, query.MemoryQueryBuilder().SetChatID(id)); err != nil {
				return err
			}

			_, err := e.Client().Rest.UpdateInteractionResponse(
				e.ApplicationID(),
				e.Token(),
				discord.NewMessageUpdateV2([]discord.LayoutComponent{
					builders.MakeSuccessContainer("해당 채팅을 삭제했어요."),
				}),
			)
			return err
		})
	})
}
