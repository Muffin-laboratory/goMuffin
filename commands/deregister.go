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
	Category: General,
	Flags:    CommandFlagsIsRegistered | CommandFlagsIsBlocked,
	Run: func(ctx *ChatInputContext) error {
		userID := ctx.Inter.User.ID

		return utils.NewMessageSender(ctx.Inter).
			AddComponents(discordgo.Container{
				Components: []discordgo.MessageComponent{
					discordgo.TextDisplay{
						Content: fmt.Sprintf("### %s 탈퇴\n- 정말로 해당 서비스에서 탈퇴하시겠어요?\n> 주의: **모든 데이터는 삭제되어요.**", ctx.Inter.Session.State.User.Username),
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								CustomID: utils.MakeDeregisterAgree(userID),
								Label:    "탈퇴",
								Style:    discordgo.DangerButton,
							},
							discordgo.Button{
								CustomID: utils.MakeDeregisterDisagree(userID),
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
	},
}

func init() {
	GetDiscommand().LoadCommand(DeregisterCommand)
}
