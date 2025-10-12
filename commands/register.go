package commands

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var RegisterCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "가입",
		Description: "이 봇에 가입해요.",
	},
	Flags: CommandFlagsIsBlocked,
	Run: func(inter *builders.InteractionCreate) error {
		userID := inter.User.ID

		if databases.GetDatabase().Users.IsUser(userID) {
			return builders.NewMessageSender(inter).
				AddComponents(builders.MakeErrorContainer("당신은 이미 가입되어있어요. 만약 탈퇴를 원하시면 /탈퇴를 이용해주세요.")).
				SetComponentsV2(true).
				SetReply(true).
				Send()
		}

		return builders.NewMessageSender(inter).
			AddComponents(
				builders.ContainerBuilder().
					AddComponents(
						builders.TextDisplayBuilder(fmt.Sprintf("### %s 가입\n해당 서비스에 가입하실려면 [개인정보처리방침](%s)과 [서비스 이용약관](%s)에 동의해야해요.",
							inter.Session.State.User.Username,
							configs.GetConfig().Service.PrivacyPolicyURL,
							configs.GetConfig().Service.TermOfServiceURL,
						)),
						builders.ActionsRowBuilder(
							builders.ButtonBuilder().
								SetStyle(discordgo.SuccessButton).
								SetLabel("동의 후 가입").
								SetCustomID(utils.MakeServiceAgree(userID)),
							builders.ButtonBuilder().
								SetStyle(discordgo.DangerButton).
								SetLabel("취소").
								SetCustomID(utils.MakeServiceDisagree(userID)),
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
	GetDiscommand().LoadCommand(RegisterCommand)
}
