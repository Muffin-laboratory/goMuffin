package chat

import (
	"context"
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func List(m any, user *discordgo.User) error {
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

	if len(data) == 0 {
		return utils.NewMessageSender(m).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "채팅이 단 하나도 없어요. 새로운 채팅을 만들거나, 대화를 시작해 채팅을 만들어주세요."})).
			SetComponentsV2(true).
			SetReply(true).
			SetEphemeral(true).
			Send()
	}

	for i, data := range data {
		sections = append(sections, discordgo.Section{
			Accessory: discordgo.Button{
				Label:    "선택",
				Style:    discordgo.SuccessButton,
				CustomID: utils.MakeSelectChat(data.Id.Hex(), i+1, user.ID),
			},
			Components: []discordgo.MessageComponent{
				discordgo.TextDisplay{
					Content: fmt.Sprintf("%d. %s\n", i+1, data.Name),
				},
			},
		})
	}

	textDisplay := discordgo.TextDisplay{Content: fmt.Sprintf("### %s님의 채팅목록", user.GlobalName)}
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
