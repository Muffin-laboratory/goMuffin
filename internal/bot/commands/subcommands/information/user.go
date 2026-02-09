package information

import (
	"context"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
)

func boolToString(k bool) string {
	if k {
		return "활성화"
	}

	return "비활성화"
}

func InfoUser(ctx context.Context, i *builders.CommandCreate) error {
	accCreatedTimestamp := i.User().ID.Time()

	dbUser, err := repository.GetDatabase().Users.FindByID(ctx, int64(i.User().ID))
	if err != nil {
		return err
	}

	currentChat, err := repository.GetDatabase().Chats.FindByID(ctx, dbUser.ChatID)
	if err != nil {
		return err
	}

	chatLength, err := repository.GetDatabase().Memory.CountDocuments(ctx, query.MemoryQueryBuilder().SetUserID(dbUser.ID))
	if err != nil {
		return err
	}

	bot, _ := i.Client().Caches.SelfUser()

	return builders.NewMessageSender(i).
		AddComponents(
			discord.NewContainer(
				discord.NewSection(
					discord.NewTextDisplayf("### %s님의 정보", *i.User().GlobalName),
					discord.NewTextDisplayf("- **디스코드 가입일**\n> %s", utils.Time(&accCreatedTimestamp, utils.RelativeTime)),
					discord.NewTextDisplayf("- **머핀봇 가입일**\n> %s", utils.Time(&dbUser.CreatedAt, utils.RelativeTime)),
				).
					WithAccessory(discord.NewThumbnail(*bot.AvatarURL())),
				discord.NewTextDisplayf("- **현재 모드**\n> `%s`", dbUser.ModeString()),
				discord.NewTextDisplayf("- **답장 멘션 사용 여부**\n> `%s`", utils.BoolToString(dbUser.ReplyUser)),
				discord.NewTextDisplayf("- **마지막 채팅 이후 12 시간이 지났을 때 새로운 채팅 생성 여부**\n> `%s`", boolToString(dbUser.CreateNewChatAfter12Hours)),
				discord.NewTextDisplayf("- **현재 채팅**\n> %s", currentChat.Name),
				discord.NewTextDisplayf("- **총 채팅량**\n> `%d`개", chatLength),
				discord.NewActionRow(
					discord.NewLinkButton("개인정보처리방침", configs.GetConfig().Service.PrivacyPolicyURL).
						WithEmoji(discord.NewComponentEmoji("🔗")),
					discord.NewLinkButton("서비스 이용약관", configs.GetConfig().Service.TermOfServiceURL).
						WithEmoji(discord.NewComponentEmoji("🔗")),
					discord.NewDangerButton("탈퇴", utils.MakeUserInformationDeregister(i.User().ID.String())),
				),
			),
		).
		SetComponentsV2(true).
		SetEphemeral(true).
		Send()
}
