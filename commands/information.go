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
		Description: "해당 봇의 정보를 알ㄹ려줘요.",
	},
	DetailedDescription: &DetailedDescription{
		Usage: fmt.Sprintf("%s정보", configs.Config.Bot.Prefix),
	},
	Category:                   General,
	RegisterApplicationCommand: true,
	RegisterMessageCommand:     true,
	Flags:                      CommandFlagsIsBlocked,
	MessageRun: func(ctx *MsgContext) error {
		return informationRun(ctx.Msg.Session, ctx.Msg)
	},
	ChatInputRun: func(ctx *ChatInputContext) error {
		return informationRun(ctx.Inter.Session, ctx.Inter)
	},
}

func informationRun(s *discordgo.Session, m any) error {
	owner, err := s.User(configs.Config.Bot.OwnerId)
	if err != nil {
		return err
	}
	return utils.NewMessageSender(m).
		AddComponents(discordgo.Container{
			Components: []discordgo.MessageComponent{
				discordgo.Section{
					Accessory: discordgo.Thumbnail{
						Media: discordgo.UnfurledMediaItem{
							URL: s.State.User.AvatarURL("512"),
						},
					},
					Components: []discordgo.MessageComponent{
						discordgo.TextDisplay{
							Content: fmt.Sprintf("### %s의 정보", s.State.User.Username),
						},
						discordgo.TextDisplay{
							Content: fmt.Sprintf("- **제작자**\n> %s", owner.Username),
						},
						discordgo.TextDisplay{
							Content: fmt.Sprintf("- **버전**\n> %s", configs.MUFFIN_VERSION),
						},
					},
				},
				discordgo.TextDisplay{
					Content: fmt.Sprintf("- **최근에 업데이트된 날짜**\n> %s", utils.Time(configs.UpdatedAt, utils.RelativeTime)),
				},
				discordgo.TextDisplay{
					Content: fmt.Sprintf("- **봇이 시작한 시각**\n> %s", utils.Time(configs.StartedAt, utils.RelativeTime)),
				},
			},
		}).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
