package chat

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/databases"
)

func SetReplyUser(i *builders.InteractionCreate, opts builders.CommandInteractionOptionsMap) error {
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

	return builders.NewMessageSender(i).
		AddComponents(builders.MakeSuccessContainer(fmt.Sprintf("대답시 답장 멘션을 성공적으로 %s했어요.", text))).
		SetComponentsV2(true).
		SetEphemeral(true).
		Send()
}
