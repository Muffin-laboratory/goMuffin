package information

import (
	"context"
	"fmt"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func InfoBot(i *utils.InteractionCreate) error {
	owner, err := i.Session.User(configs.GetConfig().Bot.OwnerID)
	if err != nil {
		return err
	}

	textLength, err := databases.GetDatabase().Texts.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return err
	}
	muffinLength, err := databases.GetDatabase().Texts.CountDocuments(context.TODO(), databases.Text{Persona: "muffin"})
	if err != nil {
		return err
	}
	learnLength, err := databases.GetDatabase().Learns.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return err
	}
	userLearnLength, err := databases.GetDatabase().Learns.CountDocuments(context.TODO(), databases.Learn{UserID: i.User.ID})
	if err != nil {
		return err
	}
	sum := textLength + learnLength

	return utils.PaginationContainerBuilder(i).
		AddContainers(
			&discordgo.Container{
				Components: []discordgo.MessageComponent{
					discordgo.Section{
						Accessory: discordgo.Thumbnail{
							Media: discordgo.UnfurledMediaItem{
								URL: i.Session.State.User.AvatarURL("512"),
							},
						},
						Components: []discordgo.MessageComponent{
							discordgo.TextDisplay{
								Content: fmt.Sprintf("### %s의 정보", i.Session.State.User.Username),
							},
							discordgo.TextDisplay{
								Content: fmt.Sprintf("- **제작자**\n> %s", owner.Username),
							},
							discordgo.TextDisplay{
								Content: fmt.Sprintf("- **버전**\n> %s", configs.MuffinVersion),
							},
						},
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **최근에 업데이트된 날짜**\n> %s", utils.Time(configs.UpdatedAt, utils.RelativeTime)),
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **봇이 시작한 시각**\n> %s", utils.Time(configs.StartedAt, utils.RelativeTime)),
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
					discordgo.Separator{},
				},
			},
			&discordgo.Container{
				Components: []discordgo.MessageComponent{
					discordgo.Section{
						Accessory: discordgo.Thumbnail{
							Media: discordgo.UnfurledMediaItem{
								URL: i.Session.State.User.AvatarURL("512"),
							},
						},
						Components: []discordgo.MessageComponent{
							discordgo.TextDisplay{
								Content: fmt.Sprintf("### 저장된 데이터량\n총합: `%d`개", sum),
							},
							discordgo.TextDisplay{
								Content: fmt.Sprintf("- **총 채팅 데이터량**\n> `%d`개", textLength),
							},
							discordgo.TextDisplay{
								Content: fmt.Sprintf("- **머핀 데이터량**\n> `%d`개", muffinLength),
							},
						},
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **총 지식 데이터량**\n> `%d`개", learnLength),
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **%s님이 가르쳐준 데이터량**\n> `%d`개", i.User.Username, userLearnLength),
					},
					discordgo.Separator{},
				},
			}).
		Start()
}
