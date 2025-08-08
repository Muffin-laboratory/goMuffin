package chat

import (
	"context"
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func List(m any, user *discordgo.User) error {
	var dbUser databases.User
	var data []databases.Chat
	var sections []discordgo.Section
	var containers []*discordgo.Container

	cur, err := databases.GetDatabase().Chats.Find(context.TODO(), databases.Chat{UserId: user.ID})
	if err != nil {
		return err
	}

	err = cur.All(context.TODO(), &data)
	if err != nil {
		return err
	}

	err = databases.GetDatabase().Users.FindOne(context.TODO(), databases.User{UserID: user.ID}).Decode(&dbUser)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return utils.NewMessageSender(m).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "채팅이 단 하나도 없어요. 새로운 채팅을 만들거나, 대화를 시작해 채팅을 만들어주세요."})).
			SetComponentsV2(true).
			SetReply(true).
			SetEphemeral(true).
			Send()
	}

	for i, data := range data {
		var isDisabled bool
		var textDisplay discordgo.TextDisplay

		if data.ID == dbUser.ChatID {
			textDisplay = discordgo.TextDisplay{
				Content: fmt.Sprintf("**%d. %s\n (선택됨)**", i+1, data.Name),
			}

			isDisabled = true
		} else {
			textDisplay = discordgo.TextDisplay{
				Content: fmt.Sprintf("%d. %s\n", i+1, data.Name),
			}

			isDisabled = false
		}

		sections = append(sections, discordgo.Section{
			Accessory: discordgo.Button{
				Label:    "선택",
				Style:    discordgo.SuccessButton,
				CustomID: utils.MakeSelectChat(data.ID.Hex(), i+1, user.ID),
				Disabled: isDisabled,
			},
			Components: []discordgo.MessageComponent{textDisplay},
		})
	}

	textDisplay := discordgo.TextDisplay{Content: fmt.Sprintf("### %s님의 채팅목록", user.GlobalName)}
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

	return utils.PaginationContainerBuilder(m).
		AddContainers(containers...).
		Start()
}
