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

var SelectChatComponent = &loader.Component{
	DeferredUpdate: true,
	Parse: func(ctx context.Context, inter *events.ComponentInteractionCreate) bool {
		customID := inter.Data.CustomID()

		if !strings.HasPrefix(customID, utils.SelectChat) {
			return false
		}

		userID := utils.GetChatUserID(customID)
		if inter.User().ID.String() != userID {
			inter.Client().Rest.UpdateInteractionResponse(
				inter.ApplicationID(),
				inter.Token(),
				discord.NewMessageUpdateBuilder().
					SetComponents(builders.MakeDeclineContainer("당신은 해당 권한이 없ㅇ어요.")).
					SetIsComponentsV2(true).
					Build(),
			)
			return false
		}
		return true
	},
	Run: func(ctx context.Context, inter *events.ComponentInteractionCreate) error {
		id, name := utils.GetChatID(inter.Data.CustomID())

		if _, err := repository.GetDatabase().Users.Update(ctx, inter.User().ID.String(), &repository.UserUpdate{
			ChatID: &id,
		}); err != nil {
			return err
		}

		_, err := inter.Client().Rest.UpdateInteractionResponse(
			inter.ApplicationID(),
			inter.Token(),
			discord.NewMessageUpdateBuilder().
				SetComponents(builders.MakeSuccessContainer("`%s`으로 채팅을 변경했어요.", name)).
				SetIsComponentsV2(true).
				Build(),
		)
		return err
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(SelectChatComponent)
}
