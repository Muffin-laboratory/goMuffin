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

var DeleteKnowledgeComponent = &loader.Component{
	DeferredUpdate: true,
	Parse: func(ctx context.Context, inter *events.ComponentInteractionCreate) bool {
		customID := inter.Data.CustomID()

		if !strings.HasPrefix(customID, utils.DeleteKnowledge) {
			return false
		}

		userID := utils.GetDeleteKnowledgeUserID(customID)
		if inter.User().ID.String() != userID {
			inter.CreateMessage(
				discord.NewMessageCreateBuilder().
					SetComponents(builders.MakeDeclineContainer("당신은 해당 권한이 없ㅇ어요.")).
					SetIsComponentsV2(true).
					SetEphemeral(true).
					Build(),
			)
			return false
		}
		return true
	},
	Run: func(ctx context.Context, inter *events.ComponentInteractionCreate) error {
		id := utils.GetDeleteKnowledgeID(inter.Data.CustomID())
		if err := repository.GetDatabase().Knowledge.DeleteByID(ctx, id); err != nil {
			return err
		}

		_, err := inter.Client().Rest.UpdateInteractionResponse(
			inter.ApplicationID(),
			inter.Token(),
			discord.NewMessageUpdateBuilder().
				SetComponents(builders.MakeSuccessContainer("해당 항목을 삭제했어요.")).
				SetIsComponentsV2(true).
				Build(),
		)

		return err
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(DeleteKnowledgeComponent)
}
