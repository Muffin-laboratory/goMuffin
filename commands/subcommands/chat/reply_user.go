package chat

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func SetReplyUser(i *utils.InteractionCreate, opts utils.CommandInteractionOptionsMap) error {
	replyUser := false
	text := "비활성화"

	if opts["활성화"].IntValue() == 1 {
		replyUser = true
		text = "활성화"
	}

	if _, err := databases.GetDatabase().Users.Update(i.User.ID, &databases.UserUpdate{
		ReplyUser: &replyUser,
	}); err != nil {
		return err
	}

	return utils.NewMessageSender(i).
		AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("대답시 답장 멘션을 성공적으로 %s했어요.", text)})).
		SetComponentsV2(true).
		SetEphemeral(true).
		Send()
}
