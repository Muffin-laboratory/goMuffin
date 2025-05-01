package commands

import (
	"fmt"
	"runtime"

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
	Category: General,
	MessageRun: func(ctx *MsgContext) {
		informationRun(ctx.Session, ctx.Msg)
	},
	ChatInputRun: func(ctx *ChatInputContext) {
		informationRun(ctx.Session, ctx.Inter)
	},
}

func informationRun(s *discordgo.Session, m any) {
	owner, _ := s.User(configs.Config.Bot.OwnerId)
	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s의 정보", s.State.User.Username),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "운영 체제",
				Value: utils.InlineCode(fmt.Sprintf("%s %s", runtime.GOARCH, runtime.GOOS)),
			},
			{
				Name:  "제작자",
				Value: utils.InlineCode(owner.Username),
			},
			{
				Name:  "버전",
				Value: utils.InlineCode(configs.MUFFIN_VERSION),
			},
			{
				Name:   "최근에 업데이트된 날짜",
				Value:  utils.Time(configs.UpdatedAt, utils.RelativeTime),
				Inline: true,
			},
			{
				Name:   "시작한 시각",
				Value:  utils.Time(configs.StartedAt, utils.RelativeTime),
				Inline: true,
			},
		},
		Color: utils.EmbedDefault,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: s.State.User.AvatarURL("512"),
		},
	}

	switch m := m.(type) {
	case *discordgo.MessageCreate:
		s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
	case *utils.InteractionCreate:
		m.Reply(&discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		})
	}
}
