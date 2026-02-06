package information

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

func InfoBot(i *builders.InteractionCreate) error {
	owner, err := i.Session.User(configs.GetConfig().Bot.OwnerID)
	if err != nil {
		return err
	}

	counts, err := repository.GetDatabase().Counts(i.Ctx, i.User.ID)
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
						AddText("### %s의 정보", i.Session.State.User.Username).
						AddText("- **제작자**\n> %s", owner.Username).
						AddText("- **버전**\n> %s", configs.MuffinVersion),
					builders.TextDisplayBuilder("- **최근에 업데이트된 날짜**\n> %s", utils.Time(configs.UpdatedAt, utils.RelativeTime)),
					builders.TextDisplayBuilder("- **봇이 시작한 시각**\n> %s", utils.Time(configs.StartedAt, utils.RelativeTime)),
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
						AddText("### 저장된 데이터 개수\n총합: `%d`개", counts.All).
						AddText("- **머핀 데이터 개수**\n> `%d`개", counts.Muffin),
					builders.TextDisplayBuilder("- **총 지식 개수**\n> `%d`개", counts.Knowledge),
					builders.TextDisplayBuilder("- **%s님이 가르쳐준 지식 개수**\n> `%d`개", i.User.GlobalName, counts.UserKnowledge),
					builders.TextDisplayBuilder("- **총 채팅방 개수**\n> `%d`개", counts.Chat),
					builders.TextDisplayBuilder("- **%s님의 채팅방 개수**\n> `%d`개", i.User.GlobalName, counts.UserChat),
					builders.TextDisplayBuilder("- **지금까지 한 채팅 개수**\n> `%d`개", counts.Memory),
					builders.TextDisplayBuilder("- **%s님이랑 지금까지 한 채팅 개수**\n> `%d`개", i.User.GlobalName, counts.UserMemory),
					builders.SeparatorBuilder(),
				),
		).
		Start()
}
