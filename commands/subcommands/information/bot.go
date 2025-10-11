package information

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func InfoBot(i *builders.InteractionCreate) error {
	owner, err := i.Session.User(configs.GetConfig().Bot.OwnerID)
	if err != nil {
		return err
	}

	counts, err := databases.GetDatabase().Counts(i.User.ID)
	if err != nil {
		return err
	}

	thumbnail := builders.ThumbnailBuilder(i.Session.State.User.AvatarURL("512"))

	return builders.PaginationContainerBuilder(i).
		AddContainers(
			builders.ContainerBuilder().
				AddComponents(
					builders.SectionBuilder().
						SetAccessory(thumbnail).
						AddText(fmt.Sprintf("### %s의 정보", i.Session.State.User.Username)).
						AddText(fmt.Sprintf("- **제작자**\n> %s", owner.Username)).
						AddText(fmt.Sprintf("- **버전**\n> %s", configs.MuffinVersion)),
					builders.TextDisplayBuilder(fmt.Sprintf("- **최근에 업데이트된 날짜**\n> %s", utils.Time(configs.UpdatedAt, utils.RelativeTime))),
					builders.TextDisplayBuilder(fmt.Sprintf("- **봇이 시작한 시각**\n> %s", utils.Time(configs.StartedAt, utils.RelativeTime))),
					builders.ActionsRowBuilder(
						builders.ButtonBuilder().
							SetStyle(discordgo.LinkButton).
							SetLabel("개인정보처리방침").
							SetURL(configs.GetConfig().Service.PrivacyPolicyURL).
							SetEmoji(discordgo.ComponentEmoji{Name: "🔗"}),
						builders.ButtonBuilder().
							SetStyle(discordgo.LinkButton).
							SetLabel("서비스 이용약관").
							SetURL(configs.GetConfig().Service.TermOfServiceURL).
							SetEmoji(discordgo.ComponentEmoji{Name: "🔗"}),
					),
					builders.SeparatorBuilder(),
				),
			builders.ContainerBuilder().
				AddComponents(
					builders.SectionBuilder().
						SetAccessory(thumbnail).
						AddText(fmt.Sprintf("### 저장된 데이터 개수\n총합: `%d`개", counts.All)).
						AddText(fmt.Sprintf("- **머핀 데이터 개수**\n> `%d`개", counts.Muffin)),
					builders.TextDisplayBuilder(fmt.Sprintf("- **총 지식 개수**\n> `%d`개", counts.Knowledge)),
					builders.TextDisplayBuilder(fmt.Sprintf("- **%s님이 가르쳐준 지식 개수**\n> `%d`개", i.User.GlobalName, counts.UserKnowledge)),
					builders.TextDisplayBuilder(fmt.Sprintf("- **총 채팅방 개수**\n> `%d`개", counts.Chat)),
					builders.TextDisplayBuilder(fmt.Sprintf("- **%s님의 채팅방 개수**\n> `%d`개", i.User.GlobalName, counts.UserChat)),
					builders.TextDisplayBuilder(fmt.Sprintf("- **지금까지 한 채팅 개수**\n> `%d`개", counts.Memory)),
					builders.TextDisplayBuilder(fmt.Sprintf("- **%s님이랑 지금까지 한 채팅 개수**\n> `%d`개", i.User.GlobalName, counts.UserMemory)),
					builders.SeparatorBuilder(),
				),
		).
		Start()
}
