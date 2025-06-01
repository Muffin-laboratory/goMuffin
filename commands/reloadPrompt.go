package commands

import (
	"log"

	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var ReloadPromptCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "프롬프트재설정",
		Description: "프롬프트를 다시 불러와요.",
	},
	Category:                   DeveloperOnly,
	RegisterApplicationCommand: false,
	MessageRun: func(ctx *MsgContext) {
		err := chatbot.ChatBot.ReloadPrompt()
		if err != nil {
			log.Fatalln(err)
			utils.NewMessageSender(ctx.Msg).
				AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "프롬프트를 다시 불러오는 데 문제가 생겼어요."})).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return
		}

		utils.NewMessageSender(ctx.Msg).
			AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: "프롬프트를 다시 불러왔어요."})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
	},
}
