package components

import (
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

var DeregisterComponent = &loader.Component{
	DeferredUpdate: true,
	Parse: func(inter *builders.InteractionCreate) bool {
		customID := inter.MessageComponentData().CustomID
		if !strings.HasPrefix(customID, utils.DeregisterAgree) && !strings.HasPrefix(customID, utils.DeregisterDisagree) {
			return false
		}

		if inter.User.ID != utils.GetDeregisterUserID(customID) {
			return false
		}
		return true
	},
	Run: func(inter *builders.InteractionCreate) error {
		customID := inter.MessageComponentData().CustomID

		switch {
		case strings.HasPrefix(customID, utils.DeregisterAgree):
			userID := inter.User.ID

			if _, err := repository.GetDatabase().Users.Delete(inter.Ctx, userID); err != nil {
				return err
			}

			if err := repository.GetDatabase().Knowledge.DeleteMany(inter.Ctx, query.KnowledgeQueryBuilder().SetUserID(userID)); err != nil {
				return err
			}

			if err := repository.GetDatabase().Memory.DeleteMany(inter.Ctx, query.MemoryQueryBuilder().SetUserID(userID)); err != nil {
				return err
			}

			return inter.EditReply(&discordgo.WebhookEdit{
				Flags: discordgo.MessageFlagsIsComponentsV2,
				Components: &[]discordgo.MessageComponent{
					builders.MakeSuccessContainer("탈퇴를 했어요.").Build(),
				},
			})
		case strings.HasPrefix(customID, utils.DeregisterDisagree):
			return inter.EditReply(&discordgo.WebhookEdit{
				Flags: discordgo.MessageFlagsIsComponentsV2,
				Components: &[]discordgo.MessageComponent{
					builders.MakeCanceledContainer("탈퇴를 거부했어요.").Build(),
				},
			})
		default:
			return nil
		}
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(DeregisterComponent)
}
