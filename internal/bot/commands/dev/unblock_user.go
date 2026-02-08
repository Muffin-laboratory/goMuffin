package dev

import (
	"context"
	"regexp"

	"github.com/LoperLee/golang-hangul-toolkit/hangul"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var UnblockCommand = &loader.Command{
	Deferred:         true,
	IsDeferEphemeral: true,
	SlashCommandCreate: &discord.SlashCommandCreate{
		Name:        "차단해제",
		Description: "유저를 차단 해제해요.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:         "유저",
				Description:  "차단 해제할 유저를 선택해요.",
				Required:     true,
				Autocomplete: true,
			},
		},
	},
	Flags: loader.CommandFlagsIsDeveloperOnlyCommand,
	Autocomplete: func(ctx context.Context, inter *events.AutocompleteInteractionCreate) error {
		var choices []discord.AutocompleteChoice
		var data []*repository.User
		var focusedValue string

		for _, opt := range inter.Data.Options {
			if opt.Focused {
				focusedValue = opt.String()
				break
			}
		}

		cur, err := repository.GetDatabase().Users.Collection.Find(ctx, bson.M{"blocked": true})
		if err != nil {
			return err
		}

		defer cur.Close(ctx)

		if err = cur.All(ctx, &data); err != nil {
			return err
		}

		for _, data := range data {
			if data.UserID == configs.GetConfig().Bot.OwnerID.String() {
				continue
			}

			user, err := inter.Client().Rest.GetUser(snowflake.MustParse(data.UserID))
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

		return inter.AutocompleteResult(choices)
	},
	Run: func(ctx context.Context, inter *builders.CommandCreate) error {
		var blocked bool
		var reason string

		userID := inter.SlashCommandInteractionData().Snowflake("유저")

		if userID == configs.GetConfig().Bot.OwnerID {
			return builders.NewMessageSender(inter).
				AddComponents(builders.MakeErrorContainer("개발자는 차단 해제를 할 수 없어요.")).
				SetComponentsV2(true).
				SetEphemeral(true).
				Send()
		}

		user, err := inter.Client().Rest.GetUser(userID)
		if err != nil {
			return err
		}

		if !repository.GetDatabase().Users.IsUser(ctx, userID.String()) {
			return builders.NewMessageSender(inter).
				AddComponents(builders.MakeErrorContainer("유저 %s은/는 해당 봇 이용자가 아니에요.", user.Username)).
				SetComponentsV2(true).
				SetEphemeral(true).
				Send()
		}

		if _, err = repository.GetDatabase().Users.Update(ctx, userID.String(), &repository.UserUpdate{
			Blocked:       &blocked,
			BlockedReason: &reason,
		}); err != nil {
			return err
		}

		return builders.NewMessageSender(inter).
			AddComponents(builders.MakeSuccessContainer("유저 %s 성공적으로 차단 해제했어요.", hangul.GetJosa(user.Username, hangul.EUL_REUL))).
			SetComponentsV2(true).
			SetEphemeral(true).
			Send()
	},
}

func init() {
	loader.GetDiscommand().LoadCommand(UnblockCommand)
}
