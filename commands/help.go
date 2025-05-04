package commands

import (
	"fmt"
	"strings"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var HelpCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "도움말",
		Description: "기본적인 사용ㅂ법이에요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "명령어",
				Description: "해당 명령어에 대ㅎ한 도움말을 볼 수 있어요.",
				Choices:     []*discordgo.ApplicationCommandOptionChoice{},
			},
		},
	},
	Aliases: []string{"도움", "명령어", "help"},
	DetailedDescription: &DetailedDescription{
		Usage:    fmt.Sprintf("%s도움말 [명령어]", configs.Config.Bot.Prefix),
		Examples: []string{fmt.Sprintf("%s도움말", configs.Config.Bot.Prefix), fmt.Sprintf("%s도움말 배워", configs.Config.Bot.Prefix)},
	},
	Category: General,
	MessageRun: func(ctx *MsgContext) {
		helpRun(ctx.Command, ctx.Session, ctx.Msg, &ctx.Args)
	},
	ChatInputRun: func(ctx *ChatInputContext) {
		helpRun(ctx.Command, ctx.Session, ctx.Inter, nil)
	},
}

func getCommandsByCategory(d *DiscommandStruct, category Category) []string {
	commands := []string{}
	for _, command := range d.Commands {
		if command.Category == category {
			commands = append(commands, fmt.Sprintf("- %s: %s", command.Name, command.Description))
		}
	}
	return commands
}

func helpRun(c *Command, s *discordgo.Session, m any, args *[]string) {
	var commandName string
	embed := &discordgo.MessageEmbed{
		Color: utils.EmbedDefault,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("버전: %s", configs.MUFFIN_VERSION),
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: s.State.User.AvatarURL("512"),
		},
	}

	switch m := m.(type) {
	case *discordgo.MessageCreate:
		commandName = Discommand.Aliases[strings.Join(*args, " ")]
	case *utils.InteractionCreate:
		if opt, ok := m.Options["명령어"]; ok {
			commandName = opt.StringValue()
		} else {
			commandName = ""
		}
	}

	if commandName == "" || Discommand.Commands[commandName] == nil {
		embed.Title = fmt.Sprintf("%s의 도움말", s.State.User.Username)
		embed.Description = utils.CodeBlock(
			"md",
			fmt.Sprintf("# 일반\n%s\n\n# 채팅\n%s",
				strings.Join(getCommandsByCategory(Discommand, General), "\n"),
				strings.Join(getCommandsByCategory(Discommand, Chatting), "\n")),
		)

		switch m := m.(type) {
		case *discordgo.MessageCreate:
			s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
		case *utils.InteractionCreate:
			m.Reply(&discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			})
		}
		return
	}

	command := Discommand.Commands[commandName]

	embed.Title = fmt.Sprintf("%s의 %s 명령어의 도움말", s.State.User.Username, command.Name)
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "설명",
			Value:  utils.InlineCode(command.Description),
			Inline: true,
		},
		{
			Name:   "사용법",
			Value:  utils.InlineCode(command.DetailedDescription.Usage),
			Inline: true,
		},
	}

	if command.Aliases != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "별칭",
			Value: utils.CodeBlock("md", strings.Join(addPrefix(command.Aliases), "\n")),
		})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "별칭",
			Value: "없음",
		})
	}

	if command.DetailedDescription.Examples != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "예시",
			Value: utils.CodeBlock("md", strings.Join(addPrefix(command.DetailedDescription.Examples), "\n")),
		})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "예시",
			Value: "없음",
		})
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
