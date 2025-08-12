package commands

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var DeregisterCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "탈퇴",
		Description: "이 봇에서 탈퇴해요.",
	},
	DetailedDescription: &DetailedDescription{
		Usage: "/탈퇴",
	},
	Category:                   General,
	RegisterMessageCommand:     true,
	RegisterApplicationCommand: true,
	Flags:                      CommandFlagsIsRegistered | CommandFlagsIsBlocked,
	MessageRun: func(ctx *MsgContext) error {
		return deregisterRun(ctx.Msg, ctx.Msg.Author.ID, ctx.Msg.Session.State.User.Username)
	},
}

func deregisterRun(m any, userId, botName string) error {
	return utils.NewMessageSender(m).
		AddComponents(discordgo.Container{
			Components: []discordgo.MessageComponent{
				discordgo.TextDisplay{
					Content: fmt.Sprintf("### %s 탈퇴\n- 정말로 해당 서비스에서 탈퇴하시겠어요?\n> 주의: **모든 데이터는 삭제되어요.**", botName),
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: utils.MakeDeregisterAgree(userId),
							Label:    "탈퇴",
							Style:    discordgo.DangerButton,
						},
						discordgo.Button{
							CustomID: utils.MakeDeregisterDisagree(userId),
							Label:    "취소",
							Style:    discordgo.PrimaryButton,
						},
					},
				},
			},
		}).
		SetComponentsV2(true).
		SetEphemeral(true).
		SetReply(true).
		Send()
}
