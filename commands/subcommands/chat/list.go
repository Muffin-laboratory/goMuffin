package chat

import (
	"context"
	"fmt"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func List(i *builders.InteractionCreate) error {
	var data []databases.Chat
	var sections []*builders.Section
	var containers []*builders.Container

	dbUser, err := databases.GetDatabase().Users.Get(i.User.ID)
	if err != nil {
		return err
	}

	if dbUser.ChattingMode == databases.ChattingMuffinMode {
		return chatSendErrorMessage(i)
	}

	cur, err := databases.GetDatabase().Chats.Find(context.TODO(), bson.M{"user_id": i.User.ID})
	if err != nil {
		return err
	}

	if err = cur.All(context.TODO(), &data); err != nil {
		return err
	}

	if len(data) == 0 {
		return builders.NewMessageSender(i).
			AddComponents(builders.MakeErrorContainer("채팅이 단 하나도 없어요. 새로운 채팅을 만들거나, 대화를 시작해 채팅을 만들어주세요.")).
			SetComponentsV2(true).
			SetReply(true).
			SetEphemeral(true).
			Send()
	}

	for _, data := range data {
		button := builders.ButtonBuilder().
			SetStyle(discordgo.SuccessButton).
			SetLabel("선택").
			SetCustomID(utils.MakeSelectChat(data.ID.Hex(), data.Name, i.User.ID))

		if data.ID == dbUser.ChatID {
			button.SetDisabled(true)

			sections = append([]*builders.Section{
				builders.SectionBuilder().
					SetAccessory(button).
					AddText(fmt.Sprintf("**%s (선택됨)**", data.Name)),
			}, sections...)

			continue
		}

		sections = append(sections,
			builders.SectionBuilder().
				SetAccessory(button).
				AddText(fmt.Sprintf("%s\n", data.Name)),
		)
	}

	textDisplay := builders.TextDisplayBuilder(fmt.Sprintf("### %s님의 채팅목록", i.User.GlobalName))
	container := builders.ContainerBuilder().AddComponents(textDisplay)
	for i, section := range sections {
		container.Components = append(container.Components, section, discordgo.Separator{})

		if (i+1)%5 == 0 {
			containers = append(containers, container)
			container = builders.ContainerBuilder().AddComponents(textDisplay)

			continue
		}
	}

	if len(container.Components) > 1 {
		containers = append(containers, container)
	}

	return builders.PaginationContainerBuilder(i).
		AddContainers(containers...).
		Start()
}
