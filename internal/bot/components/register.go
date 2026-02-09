package components

import (
	"context"
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var RegisterComponent = &loader.Component{
	DeferredUpdate: true,
	Parse: func(ctx context.Context, inter *events.ComponentInteractionCreate) bool {
		customID := inter.Data.CustomID()
		if !strings.HasPrefix(customID, utils.ServiceAgree) && !strings.HasPrefix(customID, utils.ServiceDisagree) {
			return false
		}

		if inter.User().ID.String() != utils.GetServiceUserID(customID) {
			return false
		}
		return true
	},
	Run: func(ctx context.Context, inter *events.ComponentInteractionCreate) error {
		customID := inter.Data.CustomID()

		switch {
		case strings.HasPrefix(customID, utils.ServiceAgree):
			if _, err := repository.GetDatabase().Users.Create(ctx, int64(inter.User().ID)); err != nil {
				return err
			}

			bot, _ := inter.Client().Caches.SelfUser()
			_, err := inter.Client().Rest.UpdateInteractionResponse(
				inter.ApplicationID(),
				inter.Token(),
				discord.NewMessageUpdateBuilder().
					SetComponents(builders.MakeSuccessContainer("가입을 했어요. 이제 %s의 모든 기능을 사용할 수 있어요.", bot.Username)).
					SetIsComponentsV2(true).
					Build(),
			)
			return err
		case strings.HasPrefix(customID, utils.ServiceDisagree):
			_, err := inter.Client().Rest.UpdateInteractionResponse(
				inter.ApplicationID(),
				inter.Token(),
				discord.NewMessageUpdateBuilder().
					SetComponents(builders.MakeDeclineContainer("가입을 거부했어요.")).
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
	loader.GetDiscommand().LoadComponent(RegisterComponent)
}
