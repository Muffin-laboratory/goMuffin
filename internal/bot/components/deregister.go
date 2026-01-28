package components

import (
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/commands"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

var DeregisterComponent *commands.Component = &commands.Component{
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
		err := inter.DeferUpdate()
		if err != nil {
			return err
		}

		customID := inter.MessageComponentData().CustomID
		flags := discordgo.MessageFlagsIsComponentsV2

		switch {
		case strings.HasPrefix(customID, utils.DeregisterAgree):
			userID := inter.User.ID

			if _, err := repository.GetDatabase().Users.Delete(userID); err != nil {
				return err
			}

			if _, err := repository.GetDatabase().Knowledge.DeleteByUserID(userID); err != nil {
				return err
			}

			if _, err := repository.GetDatabase().Memory.DeleteByUserID(userID); err != nil {
				return err
			}

			return inter.EditReply(&builders.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					builders.MakeSuccessContainer("탈퇴를 했어요."),
				},
			})
		case strings.HasPrefix(customID, utils.DeregisterDisagree):
			return inter.EditReply(&builders.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					builders.MakeCanceledContainer("탈퇴를 거부했어요."),
				},
			})
		default:
			return nil
		}
	},
}

func init() {
	commands.GetDiscommand().LoadComponent(DeregisterComponent)
}
