package commands

import (
	"context"
	"fmt"
	"strings"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var BlockCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "차단",
		Description: "유저를 차단해요.",
	},
	DetailedDescription: &DetailedDescription{
		Usage: fmt.Sprintf("%s차단 (유저의 ID) [사유]", configs.Config.Bot.Prefix),
	},
	Category:                   DeveloperOnly,
	RegisterApplicationCommand: false,
	RegisterMessageCommand:     true,
	Flags:                      CommandFlagsIsDeveloper,
	MessageRun: func(ctx *MsgContext) error {
		var reason string
		if len(*ctx.Args) < 1 {
			utils.NewMessageSender(ctx.Msg).
				AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "유저 ID는 필수에요."})).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return nil
		}

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

		_, err = databases.Database.Users.UpdateOne(context.TODO(),
			bson.D{{Key: "user_id", Value: userId}},
			bson.D{{
				Key: "$set",
				Value: bson.D{
					{
						Key:   "blocked",
						Value: true,
					},
					{
						Key:   "blocked_reason",
						Value: reason,
					},
				},
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
