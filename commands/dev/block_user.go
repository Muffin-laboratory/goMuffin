package dev

import (
	"fmt"
	"regexp"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"github.com/LoperLee/golang-hangul-toolkit/hangul"
	"github.com/bwmarrin/discordgo"
)

var BlockCommand = &commands.Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "차단",
		Description: "유저를 차단해요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "유저",
				Description:  "차단할 유저를 선택해요.",
				Required:     true,
				Autocomplete: true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "이유",
				Description: "해당 유저를 차단하는 이유를 적어주세요.",
			},
		},
	},
	Flags: commands.CommandFlagsIsDeveloperOnlyCommand,
	Autocomplete: func(inter *builders.InteractionCreate) error {
		var choices []*discordgo.ApplicationCommandOptionChoice
		var focusedValue string

		for _, opt := range inter.ApplicationCommandData().Options {
			if opt.Focused {
				focusedValue = opt.StringValue()
				break
			}
		}

		data, err := databases.GetDatabase().Users.All()
		if err != nil {
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
		reason := "없음"
		userID := inter.Options["유저"].StringValue()
		blocked := true

		if opt, ok := inter.Options["이유"]; ok {
			reason = opt.StringValue()
		}

		if userID == configs.GetConfig().Bot.OwnerID {
			return builders.NewMessageSender(inter).
				AddComponents(builders.MakeErrorContainer("개발자는 차단을 할 수 없어요.")).
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
			AddComponents(builders.MakeSuccessContainer(fmt.Sprintf("유저 %s 성공적으로 차단했어요.", hangul.GetJosa(user.GlobalName, hangul.EUL_REUL)))).
			SetComponentsV2(true).
			SetEphemeral(true).
			Send()
	},
}

func init() {
	commands.GetDiscommand().LoadCommand(BlockCommand)
}
