package components

import (
	"context"
	"fmt"
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

var DeleteChatComponent = &loader.Component{
	DeferredUpdate: true,
	Parse: func(inter *builders.InteractionCreate) bool {
		customID := inter.MessageComponentData().CustomID

		if !strings.HasPrefix(customID, utils.DeleteChat) && !strings.HasPrefix(customID, utils.DeleteChatCancel) {
			return false
		}

		userID := utils.GetChatUserID(customID)
		if inter.Member.User.ID != userID {
			inter.Reply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsIsComponentsV2,
				Components: []discordgo.MessageComponent{
					builders.MakeDeclineContainer("당신은 해당 권한이 없ㅇ어요."),
				},
			})
			return false
		}
		return true
	},
	Run: func(inter *builders.InteractionCreate) error {
		flags := discordgo.MessageFlagsIsComponentsV2
		customId := inter.MessageComponentData().CustomID

		if strings.HasPrefix(customId, utils.DeleteChatCancel) {
			return inter.EditReply(&builders.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					builders.MakeCanceledContainer("아무 채팅방을 삭제하지 않았어요."),
				},
			})
		}

		id, name := utils.GetChatID(inter.MessageComponentData().CustomID)

		if _, err := repository.GetDatabase().Chats.DeleteOne(context.TODO(), repository.Chat{ID: id}); err != nil {
			return err
		}

		if _, err := repository.GetDatabase().Memory.DeleteByChatID(inter.Ctx, id); err != nil {
			return err
		}

		return inter.EditReply(&builders.InteractionEdit{
			Flags: &flags,
			Components: &[]discordgo.MessageComponent{
				builders.MakeSuccessContainer(fmt.Sprintf("`%s`번을 삭제했어요.", name)),
			},
		})
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(DeleteChatComponent)
}
