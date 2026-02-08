package modals

import (
	"context"
	"strconv"
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var PaginationContainerModal = &loader.Modal{
	Parse: func(ctx context.Context, inter *events.ModalSubmitInteractionCreate) bool {
		data := inter.Data
		customID := data.CustomID

		if !strings.HasPrefix(customID, utils.PaginationContainerModal) {
			return false
		}

		id := utils.GetPaginationContainerID(customID)
		userID := utils.GetUserID(id)

		if inter.User().ID.String() != userID {
			return false
		}

		if builders.GetPaginationContainer(id) == nil {
			return false
		}

		if _, err := strconv.Atoi(data.Text(utils.PaginationContainerSetPage)); err != nil {
			inter.CreateMessage(
				discord.NewMessageCreateBuilder().
					SetComponents(builders.MakeErrorContainer("해당 값은 숫자여야해요.")).
					SetIsComponentsV2(true).
					SetEphemeral(true).
					Build(),
			)
			return false
		}

		return true
	},
	Run: func(ctx context.Context, inter *events.ModalSubmitInteractionCreate) error {
		data := inter.Data
		customID := data.CustomID
		id := utils.GetPaginationContainerID(customID)
		p := builders.GetPaginationContainer(id)

		page, _ := strconv.Atoi(data.Text(utils.PaginationContainerSetPage))

		return p.Set(inter, page)
	},
}

func init() {
	loader.GetDiscommand().LoadModal(PaginationContainerModal)
}
