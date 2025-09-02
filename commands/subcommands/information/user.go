package information

import (
	"context"
	"fmt"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func InfoUser(i *utils.InteractionCreate) error {
	accCreatedTimestamp, err := discordgo.SnowflakeTimestamp(i.User.ID)
	if err != nil {
		return err
	}

	var mufUser databases.User
	var currentChat databases.Chat

	err = databases.GetDatabase().Users.FindOne(context.TODO(), databases.User{UserID: i.User.ID}).Decode(&mufUser)
	if err != nil {
		return err
	}

	err = databases.GetDatabase().Chats.FindOne(context.TODO(), databases.Chat{ID: mufUser.ChatID}).Decode(&currentChat)
	if err != nil {
		return err
	}

	chatLength, err := databases.GetDatabase().Memory.CountDocuments(context.TODO(), databases.Memory{UserID: i.User.ID})
	if err != nil {
		return err
	}

	return utils.NewMessageSender(i).
		AddComponents(discordgo.Container{
			Components: []discordgo.MessageComponent{
				discordgo.Section{
					Accessory: discordgo.Thumbnail{Media: discordgo.UnfurledMediaItem{URL: i.User.AvatarURL("512")}},
					Components: []discordgo.MessageComponent{
						discordgo.TextDisplay{
							Content: fmt.Sprintf("### %s님의 정보", i.User.GlobalName),
						},
						discordgo.TextDisplay{
							Content: fmt.Sprintf("- **디스코드 가입일**\n> %s", utils.Time(&accCreatedTimestamp, utils.RelativeTime)),
						},
						discordgo.TextDisplay{
							Content: fmt.Sprintf("- **머핀봇 가입일**\n> %s", utils.Time(&mufUser.CreatedAt, utils.RelativeTime)),
						},
					},
				},
				discordgo.TextDisplay{
					Content: fmt.Sprintf("- **현재 모드**\n> %s", mufUser.ModeString()),
				},
				discordgo.TextDisplay{
					Content: fmt.Sprintf("- **현재 채팅**\n> %s", currentChat.Name),
				},
				discordgo.TextDisplay{
					Content: fmt.Sprintf("- **총 채팅량**\n> `%d`개", chatLength),
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "개인정보처리방침",
							URL:   configs.GetConfig().Service.PrivacyPolicyURL,
							Style: discordgo.LinkButton,
							Emoji: &discordgo.ComponentEmoji{
								Name: "🔗",
							},
						},
						discordgo.Button{
							Label: "서비스 이용약관",
							URL:   configs.GetConfig().Service.TermOfServiceURL,
							Style: discordgo.LinkButton,
							Emoji: &discordgo.ComponentEmoji{
								Name: "🔗",
							},
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: utils.MakeUserInformationDeregister(i.User.ID),
							Label:    "탈퇴",
							Style:    discordgo.DangerButton,
						},
					},
				},
			},
		}).
		SetComponentsV2(true).
		SetEphemeral(true).
		Send()
}
