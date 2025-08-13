package dev

import (
	"context"
	"fmt"
	"strings"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var BlockCommand *commands.Command = &commands.Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "차단",
		Description: "유저를 차단해요.",
	},
	DetailedDescription: commands.DetailedDescription{
		Usage: configs.AddPrefix("%s차단 (유저의 ID) [사유]"),
	},
	Category: commands.DeveloperOnly,
	Flags:    commands.CommandFlagsIsDeveloper,
	MessageRun: func(ctx *commands.MsgContext) error {
		var reason string

		userId := (*ctx.Args)[0]
		if len(*ctx.Args) >= 2 {
			reason = strings.Join((*ctx.Args)[1:], " ")
		} else {
			reason = "없음"
		}

		user, err := ctx.Msg.Session.User(userId)
		if err != nil {
			return err
		}

		_, err = databases.GetDatabase().Users.UpdateOne(context.TODO(),
			databases.User{UserID: userId},
			bson.D{{
				Key:   "$set",
				Value: databases.User{Blocked: true, BlockedReason: reason},
			}})
		if err != nil {
			return err
		}

		return utils.NewMessageSender(ctx.Msg).
			AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("- %s를 성공적으로 차단했어요.\n> 사유: %s", user.GlobalName, reason)})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
	},
}
