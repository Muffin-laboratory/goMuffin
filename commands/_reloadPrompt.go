package commands

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var ReloadPromptCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "프롬프트재설정",
		Description: "프롬프트를 다시 불러와요.",
	},
	DetailedDescription: &DetailedDescription{
		Usage: fmt.Sprintf("%s프롬프트재설정", configs.Config.Bot.Prefix),
	},
	Category:                   DeveloperOnly,
	RegisterApplicationCommand: false,
	RegisterMessageCommand:     true,
	Flags:                      CommandFlagsIsDeveloper,
	MessageRun: func(ctx *MsgContext) error {
		err := chatbot.ChatBot.ReloadPrompt()
		if err != nil {
			return err
		}

		return utils.NewMessageSender(ctx.Msg).
			AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: "프롬프트를 다시 불러왔어요."})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
	},
}
