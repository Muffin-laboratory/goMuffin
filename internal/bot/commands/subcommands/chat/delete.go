package chat

import (
	"context"
	"fmt"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

func Delete(i *builders.InteractionCreate, opts builders.CommandInteractionOptionsMap) error {
	var data []repository.Chat

	name := opts["이름"].StringValue()

	dbUser, err := repository.GetDatabase().Users.Get(i.Ctx, i.User.ID)
	if err != nil {
		return err
	}

	if dbUser.ChattingMode == repository.ChattingMuffinMode {
		return chatSendErrorMessage(i)
	}

	cur, err := repository.GetDatabase().Chats.Find(context.TODO(), repository.Chat{Name: name})
	if err != nil {
		return err
	}

	err = cur.All(context.TODO(), &data)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return builders.NewMessageSender(i).
			AddComponents(builders.MakeErrorContainer("해당하는 채팅을 찾을 수 없어요.")).
			SetComponentsV2(true).
			SetReply(true).
			Send()
	}

	if len(data) > 1 {
		var sections []*builders.Section
		var containers []*builders.Container

		for _, data := range data {
			sections = append(sections,
				builders.SectionBuilder().
					SetAccessory(
						builders.ButtonBuilder().
							SetStyle(discordgo.DangerButton).
							SetLabel("삭제").
							SetCustomID(utils.MakeDeleteChat(data.ID.Hex(), data.Name, i.User.ID)),
					).
					AddText(fmt.Sprintf("- **%s**\n", data.Name)),
			)
		}

		textDisplay := builders.TextDisplayBuilder(fmt.Sprintf("### %s님의 채팅목록\n- **주의: 이 채팅방을 삭제하면 이 채팅방의 내역을 다시는 못 써요.**", i.User.GlobalName))
		container := builders.ContainerBuilder().AddComponents(textDisplay)
		for i, section := range sections {
			container.Components = append(container.Components, section, discordgo.Separator{})

			if (i+1)%10 == 0 {
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

	return builders.NewMessageSender(i).
		AddComponents(discordgo.Container{
			Components: []discordgo.MessageComponent{
				discordgo.TextDisplay{Content: fmt.Sprintf("### 채팅 %s 삭제", name)},
				discordgo.TextDisplay{Content: "- **주의: 이 채팅방을 삭제하면 이 채팅방의 내역을 다시는 못 써요.**"},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "삭제",
							Style:    discordgo.DangerButton,
							CustomID: utils.MakeDeleteChat(data[0].ID.Hex(), name, i.User.ID),
						},
						discordgo.Button{
							Label:    "취소",
							Style:    discordgo.PrimaryButton,
							CustomID: utils.MakeDeleteChatCancel(i.User.ID),
						},
					},
				},
			},
		}).
		SetComponentsV2(true).
		SetReply(true).
		Send()

}
