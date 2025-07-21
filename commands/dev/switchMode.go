package dev

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var SwitchModeCommand *commands.Command = &commands.Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "모드전환",
		Description: "머핀봇의 대답을 변경합니다.",
	},
	DetailedDescription: &commands.DetailedDescription{
		Usage: fmt.Sprintf("%s모드전환", configs.Config.Bot.Prefix),
	},
	Category:                   commands.DeveloperOnly,
	RegisterApplicationCommand: false,
	RegisterMessageCommand:     true,
	Flags:                      commands.CommandFlagsIsDeveloper,
	MessageRun: func(ctx *commands.MsgContext) error {
		chatbot.ChatBot.SwitchMode()

		return utils.NewMessageSender(ctx.Msg).
			AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("모드를 %s로 바꾸었어요.", chatbot.ChatBot.ModeString())})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
	},
}
