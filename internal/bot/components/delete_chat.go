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

var DeleteChatComponent = &loader.Component{
	DeferredUpdate: true,
	Parse: func(ctx context.Context, inter *events.ComponentInteractionCreate) bool {
		customID := inter.Data.CustomID()

		if !strings.HasPrefix(customID, utils.DeleteChat) && !strings.HasPrefix(customID, utils.DeleteChatCancel) {
			return false
		}

		userID := utils.GetChatUserID(customID)
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
		customID := inter.Data.CustomID()

		if strings.HasPrefix(customID, utils.DeleteChatCancel) {
			_, err := inter.Client().Rest.UpdateInteractionResponse(
				inter.ApplicationID(),
				inter.Token(),
				discord.NewMessageUpdateBuilder().
					SetComponents(builders.MakeCanceledContainer("아무 채팅방을 삭제하지 않았어요.")).
					SetIsComponentsV2(true).
					Build(),
			)
			return err
		}

		id, name := utils.GetChatID(inter.Data.CustomID())

		if err := repository.GetDatabase().Chats.DeleteByID(ctx, id); err != nil {
			return err
		}

		if err := repository.GetDatabase().Memory.DeleteMany(ctx, query.MemoryQueryBuilder().SetChatID(id)); err != nil {
			return err
		}

		_, err := inter.Client().Rest.UpdateInteractionResponse(
			inter.ApplicationID(),
			inter.Token(),
			discord.NewMessageUpdateBuilder().
				SetComponents(builders.MakeSuccessContainer("`%s`번을 삭제했어요.", name)).
				SetIsComponentsV2(true).
				Build(),
		)
		return err
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(DeleteChatComponent)
}
