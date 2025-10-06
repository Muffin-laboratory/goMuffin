package chat

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func SetCreateNewChatAfter12Hours(i *utils.InteractionCreate, opts utils.CommandInteractionOptionsMap) error {
	createNewChatAfter12Hours := false
	text := "비활성화"

	if opts["활성화"].IntValue() == 1 {
		createNewChatAfter12Hours = true
		text = "활성화"
	}

	if _, err := databases.GetDatabase().Users.Update(i.User.ID, &databases.UserUpdate{
		CreateNewChatAfter12Hours: &createNewChatAfter12Hours,
	}); err != nil {
		return err
	}

	return utils.NewMessageSender(i).
		AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("마지막 대화 후 12시간이 지났을 시 새로운 채팅 생성하기를 성공적으로 %s했어요.", text)})).
		SetComponentsV2(true).
		SetEphemeral(true).
		Send()
}
