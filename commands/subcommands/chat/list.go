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

	for _, data := range data {
		button := discordgo.Button{
			Label:    "선택",
			Style:    discordgo.SuccessButton,
			CustomID: utils.MakeSelectChat(data.ID.Hex(), data.Name, i.User.ID),
		}

		if data.ID == dbUser.ChatID {
			button.Disabled = true

			sections = append([]discordgo.Section{
				{
					Accessory: button,
					Components: []discordgo.MessageComponent{
						discordgo.TextDisplay{
							Content: fmt.Sprintf("**%s (선택됨)**", data.Name),
						},
					},
				},
			}, sections...)

			continue
		}

		sections = append(sections, discordgo.Section{
			Accessory: button,
			Components: []discordgo.MessageComponent{
				discordgo.TextDisplay{
					Content: fmt.Sprintf("%s\n", data.Name),
				},
			},
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
