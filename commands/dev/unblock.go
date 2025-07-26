package dev

import (
	"context"
	"fmt"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var UnblockCommand *commands.Command = &commands.Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "차단해제",
		Description: "유저의 차단을 해제해요.",
	},
	DetailedDescription: &commands.DetailedDescription{
		Usage: configs.AddPrefix("%s차단해제 (유저의 ID)"),
	},
	Category:                   commands.DeveloperOnly,
	RegisterApplicationCommand: false,
	RegisterMessageCommand:     true,
	Flags:                      commands.CommandFlagsIsDeveloper,
	MessageRun: func(ctx *commands.MsgContext) error {
		if len(*ctx.Args) < 1 {
			utils.NewMessageSender(ctx.Msg).
				AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "유저 ID는 필수에요."})).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return nil
		}

		userId := (*ctx.Args)[0]
		user, err := ctx.Msg.Session.User(userId)
		if err != nil {
			return err
		}

		_, err = databases.GetDatabase().Users.UpdateOne(context.TODO(),
			databases.User{UserId: userId},
			bson.D{{
				Key:   "$set",
				Value: databases.User{Blocked: false, BlockedReason: ""},
			}})
		if err != nil {
			return err
		}

		return utils.NewMessageSender(ctx.Msg).
			AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("%s의 차단을 성공적으로 해제했어요.", user.GlobalName)})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
	},
}
