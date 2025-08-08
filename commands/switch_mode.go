package commands

import (
	"context"
	"fmt"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var SwitchModeCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "모드전환",
		Description: "봇의 대답 방법을 전환해요.",
	},
	DetailedDescription: &DetailedDescription{
		Usage: configs.AddPrefix("%s모드전환"),
	},
	Category:                   Chatting,
	RegisterApplicationCommand: true,
	RegisterMessageCommand:     true,
	Flags:                      CommandFlagsIsRegistered | CommandFlagsIsBlocked,
	MessageRun: func(ctx *MsgContext) error {
		return switchModeRun(ctx.Msg, ctx.Msg.Author)
	},
	ChatInputRun: func(ctx *ChatInputContext) error {
		ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return switchModeRun(ctx.Inter, ctx.Inter.User)
	},
}

func switchModeRun(m any, user *discordgo.User) error {
	var newMode databases.ChattingMode
	mode, err := databases.GetDatabase().GetUserChattingMode(user.ID)
	if err != nil {
		return err
	}

	switch mode {
	default:
		newMode = databases.ChattingAIMode
	case databases.ChattingMuffinMode:
		newMode = databases.ChattingMuffinMode
	}
	_, err = databases.GetDatabase().Users.UpdateOne(context.TODO(), databases.User{UserID: user.ID}, bson.D{{
		Key: "$set",
		Value: databases.User{
			ChattingMode: newMode,
		},
	}})
	if err != nil {
		return err
	}

	return utils.NewMessageSender(m).
		AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("모드를 성공적으로 %s로 바꿨어요.", databases.ModeString(newMode))})).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
