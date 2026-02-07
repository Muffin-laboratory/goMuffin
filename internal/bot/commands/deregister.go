package commands

import (
	"context"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
)

var DeregisterCommand = &loader.Command{
	SlashCommandCreate: &discord.SlashCommandCreate{
		Name:        "탈퇴",
		Description: "이 봇에서 탈퇴해요.",
	},
	Flags: loader.CommandFlagsIsRegistered | loader.CommandFlagsIsBlocked,
	Run: func(ctx context.Context, inter *builders.CommandCreate) error {
		userID := inter.User().ID.String()
		bot, _ := inter.Client().Caches.SelfUser()

		return builders.NewMessageSender(inter).
			AddComponents(
				discord.NewContainer(
					discord.NewTextDisplayf("### %s 탈퇴\n- 정말로 해당 서비스에서 탈퇴하시겠어요?\n> 주의: **모든 데이터는 삭제되어요.**", bot.Username)),
				discord.NewActionRow(
					discord.NewDangerButton("탈퇴", utils.MakeDeregisterAgree(userID)),
					discord.NewPrimaryButton("취소", utils.MakeDeregisterDisagree(userID)),
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
