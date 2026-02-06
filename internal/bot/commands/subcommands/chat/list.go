package chat

import (
	"fmt"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

func List(i *builders.InteractionCreate) error {
	var sections []*builders.Section
	var containers []*builders.Container

	dbUser, err := repository.GetDatabase().Users.FindByID(i.Ctx, i.User.ID)
	if err != nil {
		return err
	}

	if dbUser.ChattingMode == repository.ChattingMuffinMode {
		return chatSendErrorMessage(i)
	}

	data, err := repository.GetDatabase().Chats.Find(i.Ctx, query.ChatQueryBuilder().SetUserID(i.User.ID))
	if err != nil {
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
