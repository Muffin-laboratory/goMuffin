package commands

import (
	"fmt"

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
	DetailedDescription: DetailedDescription{
		Usage: configs.AddPrefix("%s가입"),
	},
	Category: General,
	Flags:    CommandFlagsIsBlocked,
	Run: func(ctx *ChatInputContext) error {
		return registerRun(ctx.Inter, ctx.Inter.User.ID, ctx.Inter.Session.State.User.Username)
	},
}

func registerRun(m any, userID, botName string) error {
	if databases.GetDatabase().Users.IsUser(userID) {
		utils.NewMessageSender(m).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: fmt.Sprintf("당신은 이미 가입되어있어요. 만약 탈퇴를 원하시면 %s탈퇴를 이용해주세요.", configs.GetConfig().Bot.Prefix)})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
		return nil
	}

	return utils.NewMessageSender(m).
		AddComponents(discordgo.Container{
			Components: []discordgo.MessageComponent{
				discordgo.TextDisplay{
					Content: fmt.Sprintf("### %s 가입\n해당 서비스에 가입하실려면 [개인정보처리방침](%s)과 [서비스 이용약관](%s)에 동의해야해요.",
						botName,
						configs.GetConfig().Service.PrivacyPolicyURL,
						configs.GetConfig().Service.TermOfServiceURL,
					),
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: utils.MakeServiceAgree(userID),
							Label:    "동의 후 가입",
							Style:    discordgo.SuccessButton,
						},
						discordgo.Button{
							CustomID: utils.MakeServiceDisagree(userID),
							Label:    "취소",
							Style:    discordgo.DangerButton,
						},
					},
				},
			},
		}).
		SetComponentsV2(true).
		SetReply(true).
		SetEphemeral(true).
		Send()
}
