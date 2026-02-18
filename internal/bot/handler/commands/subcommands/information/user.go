package information

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders/customid"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func boolToString(k bool) string {
	if k {
		return "활성화"
	}

	return "비활성화"
}

func InfoUser(e *handler.CommandEvent) error {
	accCreatedTimestamp := e.User().ID.Time()

	dbUser, err := repository.GetDatabase().Users.FindByID(e.Ctx, int64(e.User().ID))
	if err != nil {
		return err
	}

	currentChat, err := repository.GetDatabase().Chats.FindByID(e.Ctx, dbUser.ChatID)
	if err != nil {
		return err
	}

	chatLength, err := repository.GetDatabase().Memory.CountDocuments(e.Ctx, query.MemoryQueryBuilder().SetUserID(dbUser.ID))
	if err != nil {
		return err
	}

	bot, _ := e.Client().Caches.SelfUser()

	_, err = e.UpdateInteractionResponse(
		discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewSection(
					discord.NewTextDisplayf("### %s님의 정보", *e.User().GlobalName),
					discord.NewTextDisplayf("- **디스코드 가입일**\n> %s", builders.Time(&accCreatedTimestamp, builders.RelativeTime)),
					discord.NewTextDisplayf("- **머핀봇 가입일**\n> %s", builders.Time(&dbUser.CreatedAt, builders.RelativeTime)),
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
					discord.NewDangerButton("탈퇴", customid.MakeUserInformationDeregister(e.User().ID.String())),
				),
			),
		}),
	)
	return err
}
