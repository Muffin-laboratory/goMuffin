package chat

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func Create(m any, user *discordgo.User, name string) error {
	_, err := databases.GetDatabase().Chats.CreateChat(user.ID, name)
	if err != nil {
		return err
	}

	return utils.NewMessageSender(m).
		AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("%s를 생성했어요. 이제 현재 채팅은 %s에요.", name, name)})).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
