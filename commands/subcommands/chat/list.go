package chat

import (
	"context"
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func List(i *utils.InteractionCreate) error {
	var data []databases.Chat
	var sections []discordgo.Section
	var containers []*discordgo.Container

	dbUser, err := databases.GetDatabase().Users.Get(i.User.ID)
	if err != nil {
		return err
	}

	if dbUser.ChattingMode == databases.ChattingMuffinMode {
		return chatSendErrorMessage(i)
	}

	cur, err := databases.GetDatabase().Chats.Find(context.TODO(), databases.Chat{UserID: i.User.ID})
	if err != nil {
		return err
	}

	if err = cur.All(context.TODO(), &data); err != nil {
		return err
	}

	if len(data) == 0 {
		return utils.NewMessageSender(i).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "채팅이 단 하나도 없어요. 새로운 채팅을 만들거나, 대화를 시작해 채팅을 만들어주세요."})).
			SetComponentsV2(true).
			SetReply(true).
			SetEphemeral(true).
			Send()
	}

	// x는 원래 i인데 함수 인자의 i와 충돌나 x로 하였음.
	for x, data := range data {
		var isDisabled bool
		var textDisplay discordgo.TextDisplay

		if data.ID == dbUser.ChatID {
			textDisplay = discordgo.TextDisplay{
				Content: fmt.Sprintf("**%d. %s\n (선택됨)**", x+1, data.Name),
			}

			isDisabled = true
		} else {
			textDisplay = discordgo.TextDisplay{
				Content: fmt.Sprintf("%d. %s\n", x+1, data.Name),
			}

			isDisabled = false
		}

		sections = append(sections, discordgo.Section{
			Accessory: discordgo.Button{
				Label:    "선택",
				Style:    discordgo.SuccessButton,
				CustomID: utils.MakeSelectChat(data.ID.Hex(), x+1, i.User.ID),
				Disabled: isDisabled,
			},
			Components: []discordgo.MessageComponent{textDisplay},
		})
	}

	textDisplay := discordgo.TextDisplay{Content: fmt.Sprintf("### %s님의 채팅목록", i.User.GlobalName)}
	container := &discordgo.Container{Components: []discordgo.MessageComponent{textDisplay}}
	for i, section := range sections {
		container.Components = append(container.Components, section, discordgo.Separator{})

		if (i+1)%5 == 0 {
			containers = append(containers, container)
			container = &discordgo.Container{Components: []discordgo.MessageComponent{textDisplay}}
			continue
		}
	}

	if len(container.Components) > 1 {
		containers = append(containers, container)
	}

	return utils.PaginationContainerBuilder(i).
		AddContainers(containers...).
		Start()
}
