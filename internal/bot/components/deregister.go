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
)

var DeregisterComponent = &loader.Component{
	Middlewares: handler.Middlewares{
		middlewares.CheckIDMiddleware(),
		middlewares.TimeoutAndDeferMiddleware(loader.Timeout(), discord.InteractionTypeComponent, true, false),
	},
	Handle: func(r handler.Router) {
		r.Component(utils.DeregisterAgree, func(e *handler.ComponentEvent) error {
			userID := int64(e.User().ID)

			if _, err := repository.GetDatabase().Users.Delete(e.Ctx, userID); err != nil {
				return err
			}

			if err := repository.GetDatabase().Knowledge.DeleteMany(e.Ctx, query.KnowledgeQueryBuilder().SetUserID(userID)); err != nil {
				return err
			}

			if err := repository.GetDatabase().Memory.DeleteMany(e.Ctx, query.MemoryQueryBuilder().SetUserID(userID)); err != nil {
				return err
			}

			_, err := e.UpdateInteractionResponse(
				discord.NewMessageUpdateBuilder().
					SetComponents(builders.MakeSuccessContainer("탈퇴를 성공적으로 완료했어요.")).
					SetIsComponentsV2(true).
					Build(),
			)
			return err
		})

		r.Component(utils.DeregisterDisagree, func(e *handler.ComponentEvent) error {
			_, err := e.UpdateInteractionResponse(
				discord.NewMessageUpdateBuilder().
					SetComponents(builders.MakeSuccessContainer("탈퇴를 취소했어요.")).
					SetIsComponentsV2(true).
					Build(),
			)
			return err
		})
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(DeregisterComponent)
}
