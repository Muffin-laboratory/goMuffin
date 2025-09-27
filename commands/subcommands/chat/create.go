package chat

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func Create(i *utils.InteractionCreate, opts utils.CommandInteractionOptionsMap) error {
	name := opts["이름"].StringValue()

	dbUser, err := databases.GetDatabase().Users.Get(i.User.ID)
	if err != nil {
		return err
	}

	if dbUser.ChattingMode == databases.ChattingMuffinMode {
		return chatSendErrorMessage(i)
	}

	if _, err := databases.GetDatabase().Chats.CreateChat(i.User.ID, name); err != nil {
		return err
	}

	return utils.NewMessageSender(i).
		AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("%s를 생성했어요. 이제 현재 채팅은 %s에요.", name, name)})).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
