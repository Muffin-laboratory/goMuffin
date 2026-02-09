package chat

import (
	"context"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
)

func Delete(ctx context.Context, i *builders.CommandCreate) error {
	name := i.SlashCommandInteractionData().String("이름")

	dbUser, err := repository.GetDatabase().Users.FindByID(ctx, int64(i.User().ID))
	if err != nil {
		return err
	}

	if dbUser.ChattingMode == repository.ChattingMuffinMode {
		return chatSendErrorMessage(i)
	}

	filter := query.ChatQueryBuilder().SetUserID(dbUser.ID).SetName(name)
	data, err := repository.GetDatabase().Chats.Find(ctx, filter)
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
		var sections []discord.SectionComponent
		var containers []discord.ContainerComponent

		for _, data := range data {
			sections = append(sections,
				discord.NewSection(
					discord.NewTextDisplayf("- **%s**\n", data.Name),
				).
					WithAccessory(
						discord.NewDangerButton("삭제", utils.MakeDeleteChat(data.ID.Hex(), data.Name, i.User().ID.String())),
					),
			)
		}

		textDisplay := discord.NewTextDisplayf("### %s님의 채팅목록\n- **주의: 이 채팅방을 삭제하면 이 채팅방의 내역을 다시는 못 써요.**", *i.User().GlobalName)
		container := discord.NewContainer(textDisplay)
		for i, section := range sections {
			container = container.AddComponents(section, discord.NewSmallSeparator())

			if (i+1)%10 == 0 {
				containers = append(containers, container)
				container = container.WithComponents(textDisplay)
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
		AddComponents(
			discord.NewContainer(
				discord.NewTextDisplayf("### 채팅 %s 삭제", name),
				discord.NewTextDisplay("- **주의: 이 채팅방을 삭제하면 이 채팅방의 내역을 다시는 못 써요.**"),
				discord.NewActionRow(
					discord.NewDangerButton("삭제", utils.MakeDeleteChat(data[0].ID.Hex(), name, i.User().ID.String())),
					discord.NewPrimaryButton("취소", utils.MakeDeleteChatCancel(i.User().ID.String())),
				),
			),
		).
		SetComponentsV2(true).
		SetReply(true).
		Send()

}
