package information

import (
	"context"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
)

func InfoBot(ctx context.Context, i *builders.CommandCreate) error {
	owner, err := i.Client().Rest.GetUser(configs.GetConfig().Bot.OwnerID)
	if err != nil {
		return err
	}

	counts, err := repository.GetDatabase().Counts(ctx, i.User().ID.String())
	if err != nil {
		return err
	}

	bot, _ := i.Client().Caches.SelfUser()
	thumbnail := discord.NewThumbnail(*bot.AvatarURL())

	return builders.PaginationContainerBuilder(i).
		AddContainers(
			discord.NewContainer(
				discord.NewSection(
					discord.NewTextDisplayf("### %s의 정보", bot.Username),
					discord.NewTextDisplayf("- **제작자**\n> %s", owner.Username),
					discord.NewTextDisplayf("- **버전**\n> %s", configs.MuffinVersion),
				).
					WithAccessory(thumbnail),
				discord.NewTextDisplayf("- **최근에 업데이트된 날짜**\n> %s", utils.Time(configs.UpdatedAt, utils.RelativeTime)),
				discord.NewTextDisplayf("- **봇이 시작한 시각**\n> %s", utils.Time(configs.StartedAt, utils.RelativeTime)),
				discord.NewActionRow(
					discord.NewLinkButton("개인정보처리방침", configs.GetConfig().Service.PrivacyPolicyURL).
						WithEmoji(discord.NewComponentEmoji("🔗")),
					discord.NewLinkButton("서비스 이용약관", configs.GetConfig().Service.TermOfServiceURL).
						WithEmoji(discord.NewComponentEmoji("🔗")),
				),
				discord.NewSmallSeparator(),
			),
			discord.NewContainer(
				discord.NewSection(
					discord.NewTextDisplayf("### 저장된 데이터 개수\n총합: `%d`개", counts.All),
					discord.NewTextDisplayf("- **머핀 데이터 개수**\n> `%d`개", counts.Muffin),
				).
					WithAccessory(thumbnail),
				discord.NewTextDisplayf("- **총 지식 개수**\n> `%d`개", counts.Knowledge),
				discord.NewTextDisplayf("- **%s님이 가르쳐준 지식 개수**\n> `%d`개", *i.User().GlobalName, counts.UserKnowledge),
				discord.NewTextDisplayf("- **총 채팅방 개수**\n> `%d`개", counts.Chat),
				discord.NewTextDisplayf("- **%s님의 채팅방 개수**\n> `%d`개", *i.User().GlobalName, counts.UserChat),
				discord.NewTextDisplayf("- **지금까지 한 채팅 개수**\n> `%d`개", counts.Memory),
				discord.NewTextDisplayf("- **%s님이랑 지금까지 한 채팅 개수**\n> `%d`개", *i.User().GlobalName, counts.UserMemory),
				discord.NewSmallSeparator(),
			),
		).
		Start()
}
