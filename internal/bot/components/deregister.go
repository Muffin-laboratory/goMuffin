package components

import (
	"context"
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var DeregisterComponent = &loader.Component{
	DeferredUpdate: true,
	Parse: func(ctx context.Context, inter *events.ComponentInteractionCreate) bool {
		customID := inter.Data.CustomID()
		if !strings.HasPrefix(customID, utils.DeregisterAgree) && !strings.HasPrefix(customID, utils.DeregisterDisagree) {
			return false
		}

		if inter.User().ID.String() != utils.GetDeregisterUserID(customID) {
			return false
		}
		return true
	},
	Run: func(ctx context.Context, inter *events.ComponentInteractionCreate) error {
		customID := inter.Data.CustomID()

		switch {
		case strings.HasPrefix(customID, utils.DeregisterAgree):
			userID := inter.User().ID.String()

			if _, err := repository.GetDatabase().Users.Delete(ctx, userID); err != nil {
				return err
			}

			if err := repository.GetDatabase().Knowledge.DeleteMany(ctx, query.KnowledgeQueryBuilder().SetUserID(userID)); err != nil {
				return err
			}

			if err := repository.GetDatabase().Memory.DeleteMany(ctx, query.MemoryQueryBuilder().SetUserID(userID)); err != nil {
				return err
			}

			_, err := inter.Client().Rest.UpdateInteractionResponse(
				inter.ApplicationID(),
				inter.Token(),
				discord.NewMessageUpdateBuilder().
					SetComponents(builders.MakeSuccessContainer("탈퇴를 했어요.")).
					SetIsComponentsV2(true).
					Build(),
			)
			return err
		case strings.HasPrefix(customID, utils.DeregisterDisagree):
			_, err := inter.Client().Rest.UpdateInteractionResponse(
				inter.ApplicationID(),
				inter.Token(),
				discord.NewMessageUpdateBuilder().
					SetComponents(builders.MakeSuccessContainer("탈퇴를 거부했어요.")).
					SetIsComponentsV2(true).
					Build(),
			)
			return err
		default:
			return nil
		}
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(DeregisterComponent)
}
