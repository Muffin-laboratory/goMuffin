package commands

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders/customid"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func init() {
	const name = "가입"

	loader.GetDiscommand().RegisterCommand(discord.SlashCommandCreate{
		Name:        name,
		Description: "이 봇에 가입해요.",
	})

	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(middlewares.CheckBlocked())

		r.Command("/"+name, func(e *handler.CommandEvent) error {
			userID := e.User().ID

			if repository.GetDatabase().Users.IsUser(e.Ctx, int64(userID)) {
				return e.CreateMessage(
					discord.NewMessageCreateV2(
						builders.MakeErrorContainer("당신은 이미 가입되어있어요. 만약 탈퇴를 원하시면 `/탈퇴`를 이용해주세요."),
					).
						WithEphemeral(true),
				)
			}

			bot, _ := e.Client().Caches.SelfUser()
			return e.CreateMessage(
				discord.NewMessageCreateV2(
					discord.NewContainer(
						discord.NewTextDisplayf(
							"### %s 가입\n해당 서비스에 가입하실려면 [개인정보처리방침](%s)과 "+
								"[서비스 이용약관](%s)에 동의해야해요.",
							bot.Username,
							configs.GetConfig().Service.PrivacyPolicyURL,
							configs.GetConfig().Service.TermOfServiceURL,
						),
						discord.NewActionRow(
							discord.NewSuccessButton("동의 및 가입", customid.MakeServiceAgree(userID.String())),
							discord.NewDangerButton("취소", customid.MakeServiceDisagree(userID.String())),
						),
					),
				).
					WithEphemeral(true),
			)
		})
	})
}
