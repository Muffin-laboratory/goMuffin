package chat

import (
	"context"
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func Delete(m any, user *discordgo.User, name string) error {
	var data []databases.Chat

	cur, err := databases.GetDatabase().Chats.Find(context.TODO(), databases.Chat{Name: name})
	if err != nil {
		return err
	}

	err = cur.All(context.TODO(), &data)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return utils.NewMessageSender(m).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "해당하는 채팅을 찾을 수 없어요."})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
	}

	if len(data) > 1 {
		var sections []discordgo.Section
		var containers []*discordgo.Container

		for i, data := range data {
			sections = append(sections, discordgo.Section{
				Accessory: discordgo.Button{
					Label:    "삭제",
					Style:    discordgo.DangerButton,
					CustomID: utils.MakeDeleteChat(data.ID.Hex(), i+1, user.ID),
				},
				Components: []discordgo.MessageComponent{
					discordgo.TextDisplay{
						Content: fmt.Sprintf("%d. %s\n", i+1, data.Name),
					},
				},
			})
		}

		textDisplay := discordgo.TextDisplay{Content: fmt.Sprintf("### %s님의 채팅목록\n- **주의: 이 채팅방을 삭제하면 이 채팅방의 내역을 다시는 못 써요.**", user.GlobalName)}
		container := &discordgo.Container{Components: []discordgo.MessageComponent{textDisplay}}
		for i, section := range sections {
			container.Components = append(container.Components, section, discordgo.Separator{})

			if (i+1)%10 == 0 {
				containers = append(containers, container)
				container = &discordgo.Container{Components: []discordgo.MessageComponent{textDisplay}}
				continue
			}
		}

		if len(container.Components) > 1 {
			containers = append(containers, container)
		}

		return utils.PaginationEmbedBuilder(m).
			AddContainers(containers...).
			Start()
	}

	return utils.NewMessageSender(m).
		AddComponents(discordgo.Container{
			Components: []discordgo.MessageComponent{
				discordgo.TextDisplay{Content: fmt.Sprintf("### 채팅 %s 삭제", name)},
				discordgo.TextDisplay{Content: "- **주의: 이 채팅방을 삭제하면 이 채팅방의 내역을 다시는 못 써요.**"},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "삭제",
							Style:    discordgo.DangerButton,
							CustomID: utils.MakeDeleteChat(data[0].ID.Hex(), 0, user.ID),
						},
						discordgo.Button{
							Label:    "취소",
							Style:    discordgo.PrimaryButton,
							CustomID: utils.MakeDeleteChatCancel(user.ID),
						},
					},
				},
			},
		}).
		SetComponentsV2(true).
		SetReply(true).
		Send()

}
