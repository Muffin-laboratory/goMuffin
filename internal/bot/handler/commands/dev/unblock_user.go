package dev

import (
	"regexp"

	"github.com/LoperLee/golang-hangul-toolkit/hangul"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

func init() {
	const name = "차단해제"

	loader.GetDiscommand().RegisterDevCommand(discord.SlashCommandCreate{
		Name:        name,
		Description: "유저를 차단 해제해요.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:         "유저",
				Description:  "차단 해제할 유저를 선택해요.",
				Required:     true,
				Autocomplete: true,
			},
		},
	})

	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(
			middlewares.CheckIsDeveloper(),
			middlewares.TimeoutAndDefer(loader.Timeout(), discord.InteractionTypeApplicationCommand, false, true),
		)

		r.Autocomplete("/"+name, func(e *handler.AutocompleteEvent) error {
			var choices []discord.AutocompleteChoice

			focusedValue := e.Data.Focused().String()

			data, err := repository.GetDatabase().Users.FindBlockedUser(e.Ctx)
			if err != nil {
				return err
			}

			if len(data) > 25 {
				data = data[:25]
			}

			for _, data := range data {
				userID := snowflake.ID(data.ID)
				if userID == configs.GetConfig().Bot.OwnerID {
					continue
				}

				user, err := e.Client().Rest.GetUser(userID)
				if err != nil {
					return err
				}

				if focusedValue == "" {
					choices = append(choices, discord.AutocompleteChoiceString{
						Name:  user.Username,
						Value: user.ID.String(),
					})
				} else {
					if regexp.MustCompile(focusedValue).Match([]byte(user.Username)) {
						choices = append(choices, discord.AutocompleteChoiceString{
							Name:  user.Username,
							Value: user.ID.String(),
						})
					}
				}
			}

			return e.AutocompleteResult(choices)
		})

		r.SlashCommand("/"+name, func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
			var blocked bool
			var reason string

			userID := data.Snowflake("유저")

			if userID == configs.GetConfig().Bot.OwnerID {
				_, err := e.UpdateInteractionResponse(
					discord.NewMessageUpdateV2([]discord.LayoutComponent{
						builders.MakeErrorContainer("개발자는 차단 해제를 할 수 없어요."),
					}),
				)
				return err
			}

			user, err := e.Client().Rest.GetUser(userID)
			if err != nil {
				return err
			}

			if !repository.GetDatabase().Users.IsUser(e.Ctx, int64(userID)) {
				_, err := e.UpdateInteractionResponse(
					discord.NewMessageUpdateV2([]discord.LayoutComponent{
						builders.MakeErrorContainer("유저 %s은/는 해당 봇 이용자가 아니에요.", user.Username),
					}),
				)
				return err
			}

			if _, err = repository.GetDatabase().Users.Update(e.Ctx, int64(userID), &repository.UserUpdate{
				Blocked:       &blocked,
				BlockedReason: &reason,
			}); err != nil {
				return err
			}

			_, err = e.UpdateInteractionResponse(
				discord.NewMessageUpdateV2([]discord.LayoutComponent{
					builders.MakeSuccessContainer("유저 %s 성공적으로 차단 해제했어요.", hangul.GetJosa(user.Username, hangul.EUL_REUL)),
				}),
			)
			return err
		})
	})
}
