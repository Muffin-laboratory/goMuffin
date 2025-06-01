package commands

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var SwitchModeCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "모드전환",
		Description: "머핀봇의 대답을 변경합니다.",
	},
	Category:                   DeveloperOnly,
	RegisterApplicationCommand: false,
	MessageRun: func(ctx *MsgContext) {
		chatbot.ChatBot.SwitchMode()

		utils.NewMessageSender(ctx.Msg).
			AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("모드를 %s로 바꾸었어요.", chatbot.ChatBot.ModeString())})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
	},
}
