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

var UnblockCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "차단해제",
		Description: "유저의 차단을 해제해요.",
	},
	DetailedDescription: &DetailedDescription{
		Usage: fmt.Sprintf("%s차단해제 (유저의 ID)", configs.Config.Bot.Prefix),
	},
	Category:                   DeveloperOnly,
	RegisterApplicationCommand: false,
	RegisterMessageCommand:     true,
	Flags:                      CommandFlagsIsDeveloper,
	MessageRun: func(ctx *MsgContext) error {
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

		_, err = databases.Database.Users.UpdateOne(context.TODO(),
			bson.D{{Key: "user_id", Value: userId}},
			bson.D{{
				Key: "$set",
				Value: bson.D{
					{
						Key:   "blocked",
						Value: false,
					},
					{
						Key:   "blocked_reason",
						Value: "",
					},
				},
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
