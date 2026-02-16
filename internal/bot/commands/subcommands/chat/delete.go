package chat

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func Delete(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	name := data.String("이름")

	dbUser, err := repository.GetDatabase().Users.FindByID(e.Ctx, int64(e.User().ID))
	if err != nil {
		return err
	}

	if dbUser.ChattingMode == repository.ChattingMuffinMode {
		return chatSendErrorMessage(e)
	}

	filter := query.ChatQueryBuilder().SetUserID(dbUser.ID).SetName(name)
	chats, err := repository.GetDatabase().Chats.Find(e.Ctx, filter)
	if err != nil {
		return err
	}

	if len(chats) == 0 {
		_, err := e.UpdateInteractionResponse(
			discord.NewMessageUpdateV2([]discord.LayoutComponent{
				builders.MakeErrorContainer("해당하는 채팅을 찾을 수 없어요."),
			}),
		)
		return err
	}

	description := "- **주의: 이 채팅방을 삭제하면 이 채팅방의 내역을 다시는 못 써요.**"

	if len(chats) > 1 {
		var sections []discord.SectionComponent
		var containers []discord.ContainerComponent

		for _, data := range chats {
			sections = append(sections,
				discord.NewSection(
					discord.NewTextDisplayf("- **%s**\n", data.Name),
				).
					WithAccessory(
						discord.NewDangerButton("삭제", utils.MakeDeleteChat(data.ID.Hex(), e.User().ID.String())),
					),
			)
		}

		textDisplay := discord.NewTextDisplayf("### %s님의 채팅목록\n%s", *e.User().GlobalName, description)
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

		return builders.NewPaginatedContainer(e, true).
			AddContainers(containers...).
			Start()
	}

	_, err = e.UpdateInteractionResponse(
		discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplayf("### 채팅 %s 삭제", name),
				discord.NewTextDisplay(description),
				discord.NewActionRow(
					discord.NewDangerButton("삭제", utils.MakeDeleteChat(chats[0].ID.Hex(), e.User().ID.String())),
					discord.NewPrimaryButton("취소", utils.MakeDeleteChatCancel(e.User().ID.String())),
				),
			),
		}),
	)
	return err
}
