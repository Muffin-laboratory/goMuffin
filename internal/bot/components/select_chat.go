package components

import (
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

var SelectChatComponent = &loader.Component{
	DeferredUpdate: true,
	Parse: func(inter *builders.InteractionCreate) bool {
		customID := inter.MessageComponentData().CustomID

		if !strings.HasPrefix(customID, utils.SelectChat) {
			return false
		}

		userID := utils.GetChatUserID(customID)
		if inter.User.ID != userID {
			inter.Reply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsIsComponentsV2,
				Components: []discordgo.MessageComponent{
					builders.MakeDeclineContainer("당신은 해당 권한이 없ㅇ어요.").Build(),
				},
			})
			return false
		}
		return true
	},
	Run: func(inter *builders.InteractionCreate) error {
		id, name := utils.GetChatID(inter.MessageComponentData().CustomID)

		if _, err := repository.GetDatabase().Users.Update(inter.Ctx, inter.User.ID, &repository.UserUpdate{
			ChatID: &id,
		}); err != nil {
			return err
		}

		return inter.EditReply(&discordgo.WebhookEdit{
			Flags: discordgo.MessageFlagsIsComponentsV2,
			Components: &[]discordgo.MessageComponent{
				builders.MakeSuccessContainer("`%s`으로 채팅을 변경했어요.", name).Build(),
			},
		})
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(SelectChatComponent)
}
