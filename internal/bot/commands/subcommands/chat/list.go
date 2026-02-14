package chat

import (
	"context"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
)

func List(ctx context.Context, i *builders.CommandCreate) error {
	var sections []discord.SectionComponent
	var containers []discord.ContainerComponent

	dbUser, err := repository.GetDatabase().Users.FindByID(ctx, int64(i.User().ID))
	if err != nil {
		return err
	}

	if dbUser.ChattingMode == repository.ChattingMuffinMode {
		return chatSendErrorMessage(i)
	}

	data, err := repository.GetDatabase().Chats.Find(ctx, query.ChatQueryBuilder().SetUserID(int64(i.User().ID)))
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
		button := discord.NewSuccessButton("선택", utils.MakeSelectChat(data.ID.Hex(), i.User().ID.String()))

		if data.ID == dbUser.ChatID {
			button = button.WithDisabled(true)

			sections = append([]discord.SectionComponent{
				discord.NewSection(
					discord.NewTextDisplayf("**%s (선택됨)**", data.Name),
				).
					WithAccessory(button),
			}, sections...)

			continue
		}

		sections = append(sections,
			discord.NewSection(
				discord.NewTextDisplayf("%s\n", data.Name),
			).
				WithAccessory(button),
		)
	}

	textDisplay := discord.NewTextDisplayf("### %s님의 채팅목록", *i.User().GlobalName)
	container := discord.NewContainer(textDisplay)
	for i, section := range sections {
		container = container.AddComponents(section, discord.NewSmallSeparator())

		if (i+1)%5 == 0 {
			containers = append(containers, container)
			container = container.WithComponents(textDisplay)

			continue
		}
	}

	if len(container.Components) > 1 {
		containers = append(containers, container)
	}

	return builders.PaginationContainerBuilder(i, true).
		AddContainers(containers...).
		Start()
}
