package commands

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var InformationCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "정보",
		Description: "해당 봇의 정보를 알려줘요.",
	},
	DetailedDescription: DetailedDescription{
		Usage: "/정보",
	},
	Category: General,
	Flags:    CommandFlagsIsBlocked,
	Run: func(ctx *ChatInputContext) error {
		owner, err := ctx.Inter.Session.User(configs.GetConfig().Bot.OwnerID)
		if err != nil {
			return err
		}

		return utils.NewMessageSender(ctx.Inter).
			AddComponents(discordgo.Container{
				Components: []discordgo.MessageComponent{
					discordgo.Section{
						Accessory: discordgo.Thumbnail{
							Media: discordgo.UnfurledMediaItem{
								URL: ctx.Inter.Session.State.User.AvatarURL("512"),
							},
						},
						Components: []discordgo.MessageComponent{
							discordgo.TextDisplay{
								Content: fmt.Sprintf("### %s의 정보", ctx.Inter.Session.State.User.Username),
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
				},
			}).
			SetComponentsV2(true).
			SetReply(true).
			Send()

	},
}
