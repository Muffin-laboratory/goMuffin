package commands

import (
	"fmt"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

var DeregisterCommand = &loader.Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "탈퇴",
		Description: "이 봇에서 탈퇴해요.",
	},
	Flags: loader.CommandFlagsIsRegistered | loader.CommandFlagsIsBlocked,
	Run: func(inter *builders.InteractionCreate) error {
		userID := inter.User.ID

		return builders.NewMessageSender(inter).
			AddComponents(
				builders.ContainerBuilder().
					AddComponents(
						builders.TextDisplayBuilder(fmt.Sprintf("### %s 탈퇴\n- 정말로 해당 서비스에서 탈퇴하시겠어요?\n> 주의: **모든 데이터는 삭제되어요.**", inter.Session.State.User.Username)),
						builders.ActionsRowBuilder(
							builders.ButtonBuilder().
								SetStyle(discordgo.DangerButton).
								SetLabel("탈퇴").
								SetCustomID(utils.MakeDeregisterAgree(userID)),
							builders.ButtonBuilder().
								SetStyle(discordgo.PrimaryButton).
								SetLabel("취소").
								SetCustomID(utils.MakeDeregisterDisagree(userID)),
						),
					),
			).
			SetComponentsV2(true).
			SetEphemeral(true).
			SetReply(true).
			Send()
	},
}

func init() {
	loader.GetDiscommand().LoadCommand(DeregisterCommand)
}
