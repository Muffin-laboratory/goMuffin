package components

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var DeleteChatComponent = &loader.Component{
	Middlewares: handler.Middlewares{
		middlewares.TimeoutAndDeferMiddleware(loader.Timeout(), discord.InteractionTypeComponent, true, false),
	},
	Handle: func(r handler.Router) {
		r.Component(utils.DeleteChat+"/{value}/{user_id}", func(e *handler.ComponentEvent) error {
			value := e.Vars["value"]

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

			if value == "cancel" {
				_, err := e.UpdateInteractionResponse(
					discord.NewMessageUpdateBuilder().
						SetComponents(builders.MakeCanceledContainer("아무 채팅방을 삭제하지 않았어요.")).
						SetIsComponentsV2(true).
						Build(),
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
				discord.NewMessageUpdateBuilder().
					SetComponents(builders.MakeSuccessContainer("해당 채팅을 삭제했어요.")).
					SetIsComponentsV2(true).
					Build(),
			)
			return err
		})
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(DeleteChatComponent)
}
