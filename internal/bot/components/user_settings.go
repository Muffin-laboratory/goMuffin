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

		return builders.GetUserSettings(id) != nil
	},
	Run: func(ctx context.Context, inter *events.ComponentInteractionCreate) error {
		customID := inter.Data.CustomID()
		settings := builders.GetUserSettings(utils.GetUserSettingsID(customID))

		switch {
		case strings.HasPrefix(customID, utils.UserSettingsChattingMode):
			newMode := repository.ChattingAIMode

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
				discord.NewMessageUpdateV2([]discord.LayoutComponent{
					builders.MakeSuccessContainer("- 봇의 대화 설정을 성공적으로 바꾸었어요."),
				}),
			)
			return err
		default:
			return nil
		}

	ReturnSettings:
		_, err := inter.Client().Rest.UpdateInteractionResponse(
			inter.ApplicationID(),
			inter.Token(),
			discord.NewMessageUpdateV2([]discord.LayoutComponent{
				settings.MakeContainer(),
			}),
		)
		return err
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(UserSettingsComponent)
}
