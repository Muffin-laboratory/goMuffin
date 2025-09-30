package information

import (
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

	counts, err := databases.GetDatabase().Counts(i.User.ID)
	if err != nil {
		return err
	}

	thumbnail := &discordgo.Thumbnail{
		Media: discordgo.UnfurledMediaItem{
			URL: i.Session.State.User.AvatarURL("512"),
		},
	}

	return utils.PaginationContainerBuilder(i).
		AddContainers(
			&discordgo.Container{
				Components: []discordgo.MessageComponent{
					discordgo.Section{
						Accessory: thumbnail,
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
						Accessory: thumbnail,
						Components: []discordgo.MessageComponent{
							discordgo.TextDisplay{
								Content: fmt.Sprintf("### 저장된 데이터 개수\n총합: `%d`개", counts.All),
							},
							discordgo.TextDisplay{
								Content: fmt.Sprintf("- **머핀 데이터 개수**\n> `%d`개", counts.Muffin),
							},
						},
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **총 지식 개수**\n> `%d`개", counts.Knowledge),
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **%s님이 가르쳐준 지식 개수**\n> `%d`개", i.User.GlobalName, counts.UserKnowledge),
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **총 채팅방 개수**\n> `%d`개", counts.Chat),
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **%s님의 채팅방 개수**\n> `%d`개", i.User.GlobalName, counts.UserChat),
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **지금까지 한 채팅 개수**\n> `%d`개", counts.Memory),
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **%s님이랑 지금까지 한 채팅 개수**\n> `%d`개", i.User.GlobalName, counts.UserMemory),
					},
					discordgo.Separator{},
				},
			}).
		Start()
}
