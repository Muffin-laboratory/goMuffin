package chat

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

func Delete(i *builders.InteractionCreate, opts *discordgo.ApplicationCommandInteractionDataOption) error {
	name := opts.GetOption("이름").StringValue()

	dbUser, err := repository.GetDatabase().Users.FindByID(i.Ctx, i.User.ID)
	if err != nil {
		return err
	}

	if dbUser.ChattingMode == repository.ChattingMuffinMode {
		return chatSendErrorMessage(i)
	}

	filter := query.ChatQueryBuilder().SetUserID(dbUser.UserID).SetName(name)
	data, err := repository.GetDatabase().Chats.Find(i.Ctx, filter)
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
					AddText("- **%s**\n", data.Name),
			)
		}

		textDisplay := builders.TextDisplayBuilder("### %s님의 채팅목록\n- **주의: 이 채팅방을 삭제하면 이 채팅방의 내역을 다시는 못 써요.**", i.User.GlobalName)
		container := builders.ContainerBuilder().AddComponents(textDisplay)
		for i, section := range sections {
			container.AddComponents(section, builders.SeparatorBuilder())

			if (i+1)%10 == 0 {
				containers = append(containers, container)
				container = builders.ContainerBuilder().AddComponents(textDisplay)
				continue
			}
		}

		if container.GetComponentsLength() > 1 {
			containers = append(containers, container)
		}

		return builders.PaginationContainerBuilder(i).
			AddContainers(containers...).
			Start()
	}

	return builders.NewMessageSender(i).
		AddComponents(
			builders.ContainerBuilder().
				AddText("### 채팅 %s 삭제", name).
				AddText("- **주의: 이 채팅방을 삭제하면 이 채팅방의 내역을 다시는 못 써요.**").
				AddComponents(
					builders.ActionsRowBuilder(
						builders.ButtonBuilder().
							SetLabel("삭제").
							SetStyle(discordgo.DangerButton).
							SetCustomID(utils.MakeDeleteChat(data[0].ID.Hex(), name, i.User.ID)),
						builders.ButtonBuilder().
							SetLabel("취소").
							SetStyle(discordgo.PrimaryButton).
							SetCustomID(utils.MakeDeleteChatCancel(i.User.ID)),
					),
				),
		).
		SetComponentsV2(true).
		SetReply(true).
		Send()

}
