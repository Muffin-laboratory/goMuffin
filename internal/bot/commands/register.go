package commands

import (
	"context"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
)

var RegisterCommand = &loader.Command{
	Deferred:         true,
	IsDeferEphemeral: true,
	SlashCommandCreate: &discord.SlashCommandCreate{
		Name:        "가입",
		Description: "이 봇에 가입해요.",
	},
	Flags: loader.CommandFlagsIsBlocked,
	Run: func(ctx context.Context, inter *builders.CommandCreate) error {
		userID := inter.User().ID.String()
		bot, _ := inter.Client().Caches.SelfUser()

		if repository.GetDatabase().Users.IsUser(ctx, userID) {
			return builders.NewMessageSender(inter).
				AddComponents(builders.MakeErrorContainer("당신은 이미 가입되어있어요. 만약 탈퇴를 원하시면 /탈퇴를 이용해주세요.")).
				SetComponentsV2(true).
				SetReply(true).
				Send()
		}

		return builders.NewMessageSender(inter).
			AddComponents(
				discord.NewContainer(
					discord.NewTextDisplayf("### %s 가입\n해당 서비스에 가입하실려면 [개인정보처리방침](%s)과 [서비스 이용약관](%s)에 동의해야해요.",
						bot.Username,
						configs.GetConfig().Service.PrivacyPolicyURL,
						configs.GetConfig().Service.TermOfServiceURL,
					),
					discord.NewActionRow(
						discord.NewSuccessButton("동의 및 가입", utils.MakeServiceAgree(userID)),
						discord.NewDangerButton("취소", utils.MakeServiceDisagree(userID)),
					),
				),
			).
			SetComponentsV2(true).
			SetReply(true).
			SetEphemeral(true).
			Send()
	},
}

func init() {
	loader.GetDiscommand().LoadCommand(RegisterCommand)
}
