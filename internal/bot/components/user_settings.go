package components

import (
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/commands"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

var UserSettingsComponent = &commands.Component{
	DeferredUpdate: true,
	Parse: func(inter *builders.InteractionCreate) bool {
		customID := inter.MessageComponentData().CustomID

		isChattingMode := strings.HasPrefix(customID, utils.UserSettingsChattingMode)
		isReplyUser := strings.HasPrefix(customID, utils.UserSettingsReplyUser)
		is12Hours := strings.HasPrefix(customID, utils.UserSettings12Hours)
		isSubmit := strings.HasPrefix(customID, utils.UserSettingsSubmit)

		if !isChattingMode && !isReplyUser && !is12Hours && !isSubmit {
			return false
		}

		id := utils.GetUserSettingsID(customID)
		if inter.User.ID != utils.GetUserID(id) {
			return false
		}

		return repository.GetUserSettings(id) != nil
	},
	Run: func(inter *builders.InteractionCreate) error {
		flags := discordgo.MessageFlagsIsComponentsV2
		customID := inter.MessageComponentData().CustomID
		settings := repository.GetUserSettings(utils.GetUserSettingsID(customID))

		switch {
		case strings.HasPrefix(customID, utils.UserSettingsChattingMode):
			var newMode repository.ChattingMode

			if settings.ChattingMode == repository.ChattingAIMode {
				newMode = repository.ChattingMuffinMode
			}

			settings.ChattingMode = newMode

			return inter.EditReply(&builders.InteractionEdit{
				Flags:      &flags,
				Components: &[]discordgo.MessageComponent{settings.MakeContainer().Build()},
			})
		case strings.HasPrefix(customID, utils.UserSettingsReplyUser):
			settings.ReplyUser = !settings.ReplyUser

			return inter.EditReply(&builders.InteractionEdit{
				Flags:      &flags,
				Components: &[]discordgo.MessageComponent{settings.MakeContainer().Build()},
			})
		case strings.HasPrefix(customID, utils.UserSettings12Hours):
			settings.CreateNewChatAfter12Hours = !settings.CreateNewChatAfter12Hours

			return inter.EditReply(&builders.InteractionEdit{
				Flags:      &flags,
				Components: &[]discordgo.MessageComponent{settings.MakeContainer().Build()},
			})
		case strings.HasPrefix(customID, utils.UserSettingsSubmit):
			if err := settings.Submit(inter.Ctx); err != nil {
				return err
			}

			return inter.EditReply(&builders.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					builders.ContainerBuilder().
						AddComponents(
							builders.SectionBuilder().
								SetAccessory(builders.ThumbnailBuilder(inter.User.AvatarURL("512"))).
								AddText("### 채팅 설정\n- 봇의 설정을 성공적으로 바꾸었어요."),
						),
				},
			})
		default:
			return nil
		}
	},
}

func init() {
	commands.GetDiscommand().LoadComponent(UserSettingsComponent)
}
