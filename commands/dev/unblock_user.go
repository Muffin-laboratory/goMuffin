package dev

import (
	"context"
	"fmt"
	"regexp"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"github.com/LoperLee/golang-hangul-toolkit/hangul"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var UnblockCommand = &commands.Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "차단해제",
		Description: "유저를 차단 해제해요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "유저",
				Description:  "차단 해제할 유저를 선택해요.",
				Required:     true,
				Autocomplete: true,
			},
		},
	},
	Flags: commands.CommandFlagsIsDeveloperOnlyCommand,
	Autocomplete: func(inter *builders.InteractionCreate) error {
		var choices []*discordgo.ApplicationCommandOptionChoice
		var data []*databases.User
		var focusedValue string

		for _, opt := range inter.ApplicationCommandData().Options {
			if opt.Focused {
				focusedValue = opt.StringValue()
				break
			}
		}

		cur, err := databases.GetDatabase().Users.Collection.Find(context.TODO(), bson.M{"blocked": true})
		if err != nil {
			return err
		}

		defer cur.Close(context.TODO())

		if err = cur.All(context.TODO(), &data); err != nil {
			return err
		}

		for _, data := range data {
			if data.UserID == configs.GetConfig().Bot.OwnerID {
				continue
			}

			user, err := inter.Session.User(data.UserID)
			if err != nil {
				return err
			}

			if focusedValue == "" {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  user.GlobalName,
					Value: user.ID,
				})
			} else {
				if regexp.MustCompile(focusedValue).Match([]byte(user.GlobalName)) {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  user.GlobalName,
						Value: user.ID,
					})
				}
			}
		}

		return inter.Autocomplete(choices)
	},
	Run: func(inter *builders.InteractionCreate) error {
		var blocked bool
		var reason string

		userID := inter.Options["유저"].StringValue()

		if userID == configs.GetConfig().Bot.OwnerID {
			return builders.NewMessageSender(inter).
				AddComponents(builders.MakeErrorContainer("개발자는 차단 해제를 할 수 없어요.")).
				SetComponentsV2(true).
				SetEphemeral(true).
				Send()
		}

		user, err := inter.Session.User(userID)
		if err != nil {
			return err
		}

		if !databases.GetDatabase().Users.IsUser(userID) {
			return builders.NewMessageSender(inter).
				AddComponents(builders.MakeErrorContainer(fmt.Sprintf("유저 %s은/는 해당 봇 이용자가 아니에요.", user.GlobalName))).
				SetComponentsV2(true).
				SetEphemeral(true).
				Send()
		}

		if _, err = databases.GetDatabase().Users.Update(userID, &databases.UserUpdate{
			Blocked:       &blocked,
			BlockedReason: &reason,
		}); err != nil {
			return err
		}

		return builders.NewMessageSender(inter).
			AddComponents(builders.MakeSuccessContainer(fmt.Sprintf("유저 %s 성공적으로 차단 해제했어요.", hangul.GetJosa(user.GlobalName, hangul.EUL_REUL)))).
			SetComponentsV2(true).
			SetEphemeral(true).
			Send()
	},
}

func init() {
	commands.GetDiscommand().LoadCommand(UnblockCommand)
}
