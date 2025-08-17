package chat

import (
	"context"
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func SwitchMode(i *utils.InteractionCreate) error {
	var newMode databases.ChattingMode

	mode, err := databases.GetDatabase().Users.GetUserChattingMode(i.User.ID)
	if err != nil {
		return err
	}

	switch mode {
	case databases.ChattingMuffinMode:
		newMode = databases.ChattingAIMode
	case databases.ChattingAIMode:
		newMode = databases.ChattingMuffinMode
	}

	_, err = databases.GetDatabase().Users.UpdateOne(context.TODO(), databases.User{UserID: i.User.ID}, bson.D{{
		Key: "$set",
		Value: databases.User{
			ChattingMode: databases.ChattingAIMode,
		},
	}})
	if err != nil {
		return err
	}

	return utils.NewMessageSender(i).
		AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("모드를 성공적으로 %s로 바꿨어요.", databases.ModeString(newMode))})).
		SetComponentsV2(true).
		Send()
}
