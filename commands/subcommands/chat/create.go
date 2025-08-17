package chat

import (
	"context"
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func Create(i *utils.InteractionCreate, name string) error {
	var dbUser databases.User

	err := databases.GetDatabase().Users.FindOne(context.TODO(), databases.User{UserID: i.User.ID}).Decode(&dbUser)
	if err != nil {
		return err
	}

	if dbUser.ChattingMode == databases.ChattingMuffinMode {
		return chatSendErrorMessage(i)
	}

	_, err = databases.GetDatabase().Chats.CreateChat(i.User.ID, name)
	if err != nil {
		return err
	}

	return utils.NewMessageSender(i).
		AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("%s를 생성했어요. 이제 현재 채팅은 %s에요.", name, name)})).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
