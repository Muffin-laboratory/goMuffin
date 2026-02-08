package components

import (
	"context"
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var UserSettingsComponent = &loader.Component{
	DeferredUpdate: true,
	Parse: func(ctx context.Context, inter *events.ComponentInteractionCreate) bool {
		customID := inter.Data.CustomID()

		isChattingMode := strings.HasPrefix(customID, utils.UserSettingsChattingMode)
		isReplyUser := strings.HasPrefix(customID, utils.UserSettingsReplyUser)
		is12Hours := strings.HasPrefix(customID, utils.UserSettings12Hours)
		isSubmit := strings.HasPrefix(customID, utils.UserSettingsSubmit)

		if !isChattingMode && !isReplyUser && !is12Hours && !isSubmit {
			return false
		}

		id := utils.GetUserSettingsID(customID)
		if inter.User().ID.String() != utils.GetUserID(id) {
			return false
		}

		return repository.GetUserSettings(id) != nil
	},
	Run: func(ctx context.Context, inter *events.ComponentInteractionCreate) error {
		customID := inter.Data.CustomID()
		settings := repository.GetUserSettings(utils.GetUserSettingsID(customID))

		switch {
		case strings.HasPrefix(customID, utils.UserSettingsChattingMode):
			var newMode repository.ChattingMode

			if settings.ChattingMode == repository.ChattingAIMode {
				newMode = repository.ChattingMuffinMode
			}

			settings.ChattingMode = newMode

			goto ReturnSettings
		case strings.HasPrefix(customID, utils.UserSettingsReplyUser):
			settings.ReplyUser = !settings.ReplyUser

			goto ReturnSettings
		case strings.HasPrefix(customID, utils.UserSettings12Hours):
			settings.CreateNewChatAfter12Hours = !settings.CreateNewChatAfter12Hours

			goto ReturnSettings
		case strings.HasPrefix(customID, utils.UserSettingsSubmit):
			if err := settings.Submit(ctx); err != nil {
				return err
			}

			_, err := inter.Client().Rest.UpdateInteractionResponse(
				inter.ApplicationID(),
				inter.Token(),
				discord.NewMessageUpdateBuilder().
					SetComponents(
						discord.NewContainer(
							discord.NewSection(
								discord.NewTextDisplay("### 채팅 설정\n- 봇의 설정을 성공적으로 바꾸었어요."),
							).
								WithAccessory(
									discord.NewThumbnail(*inter.User().AvatarURL()),
								),
						),
					).
					SetIsComponentsV2(true).
					Build(),
			)
			return err
		default:
			return nil
		}

	ReturnSettings:
		_, err := inter.Client().Rest.UpdateInteractionResponse(
			inter.ApplicationID(),
			inter.Token(),
			discord.NewMessageUpdateBuilder().
				SetComponents(settings.MakeContainer()).
				SetIsComponentsV2(true).
				Build(),
		)
		return err
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(UserSettingsComponent)
}
