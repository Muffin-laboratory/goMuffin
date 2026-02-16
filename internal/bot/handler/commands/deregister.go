package commands

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

func init() {
	const name = "탈퇴"

	loader.GetDiscommand().RegisterCommand(discord.SlashCommandCreate{
		Name:        name,
		Description: "이 봇에서 탈퇴해요.",
	})

	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(middlewares.CheckUserAndBlockedMiddleware())

		r.Command("/"+name, func(e *handler.CommandEvent) error {
			return HandleDeregister(e.ApplicationCommandInteractionCreate)
		})
	})
}

func HandleDeregister(inter *events.ApplicationCommandInteractionCreate) error {
	userID := inter.User().ID.String()
	bot, _ := inter.Client().Caches.SelfUser()

	return inter.CreateMessage(
		discord.NewMessageCreateV2(
			discord.NewContainer(
				discord.NewTextDisplayf("### %s 탈퇴\n- 정말로 해당 서비스에서 탈퇴하시겠어요?\n> 주의: **모든 데이터는 삭제되어요.**", bot.Username),
				discord.NewActionRow(
					discord.NewDangerButton("탈퇴", utils.MakeDeregisterAgree(userID)),
					discord.NewPrimaryButton("취소", utils.MakeDeregisterDisagree(userID)),
				),
			),
		).
			WithEphemeral(true),
	)
}
