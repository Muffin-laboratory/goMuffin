package information

import (
	"context"
	"fmt"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/repository"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func boolToString(k bool) string {
	if k {
		return "활성화"
	}

	return "비활성화"
}

func InfoUser(i *builders.InteractionCreate) error {
	accCreatedTimestamp, err := discordgo.SnowflakeTimestamp(i.User.ID)
	if err != nil {
		return err
	}

	var currentChat repository.Chat

	dbUser, err := repository.GetDatabase().Users.Get(i.User.ID)
	if err != nil {
		return err
	}

	if err = repository.GetDatabase().Chats.FindOne(context.TODO(), repository.Chat{
		ID: dbUser.ChatID,
	}).Decode(&currentChat); err != nil {
		return err
	}

	chatLength, err := repository.GetDatabase().Memory.Collection.CountDocuments(context.TODO(), repository.Memory{UserID: i.User.ID})
	if err != nil {
		return err
	}

	return builders.NewMessageSender(i).
		AddComponents(
			builders.ContainerBuilder().
				AddComponents(
					builders.SectionBuilder().
						SetAccessory(builders.ThumbnailBuilder(i.User.AvatarURL("512"))).
						AddText(fmt.Sprintf("### %s님의 정보", i.User.GlobalName)).
						AddText(fmt.Sprintf("- **디스코드 가입일**\n> %s", utils.Time(&accCreatedTimestamp, utils.RelativeTime))).
						AddText(fmt.Sprintf("- **머핀봇 가입일**\n> %s", utils.Time(&dbUser.CreatedAt, utils.RelativeTime))),
					builders.TextDisplayBuilder(fmt.Sprintf("- **현재 모드**\n> `%s`", dbUser.ModeString())),
					builders.TextDisplayBuilder(fmt.Sprintf("- **답장 멘션 사용 여부**\n> `%s`", boolToString(dbUser.ReplyUser))),
					builders.TextDisplayBuilder(fmt.Sprintf("- **마지막 채팅 이후 12 시간이 지났을 때 새로운 채팅 생성 여부**\n> `%s`", boolToString(dbUser.CreateNewChatAfter12Hours))),
					builders.TextDisplayBuilder(fmt.Sprintf("- **현재 채팅**\n> %s", currentChat.Name)),
					builders.TextDisplayBuilder(fmt.Sprintf("- **총 채팅량**\n> `%d`개", chatLength)),
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
						builders.ButtonBuilder().
							SetStyle(discordgo.DangerButton).
							SetLabel("탈퇴").
							SetCustomID(utils.MakeUserInformationDeregister(i.User.ID)),
					),
				),
		).
		SetComponentsV2(true).
		SetEphemeral(true).
		Send()
}
