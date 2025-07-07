package commands

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var SwitchModeCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "모드전환",
		Description: "머핀봇의 대답을 변경합니다.",
	},
	DetailedDescription: &DetailedDescription{
		Usage: fmt.Sprintf("%s모드전환", configs.Config.Bot.Prefix),
	},
	Category:                   DeveloperOnly,
	RegisterApplicationCommand: false,
	RegisterMessageCommand:     true,
	Flags:                      CommandFlagsIsDeveloper,
	MessageRun: func(ctx *MsgContext) error {
		chatbot.ChatBot.SwitchMode()

		return utils.NewMessageSender(ctx.Msg).
			AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("모드를 %s로 바꾸었어요.", chatbot.ChatBot.ModeString())})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
	},
}
