package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/LoperLee/golang-hangul-toolkit/hangul"
	"github.com/bwmarrin/discordgo"
)

var LearnCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "배워",
		Description: "단어를 가르치는 명령ㅇ어에요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "단어",
				Description: "등록할 단어를 입력해주세요.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "대답",
				Description: "해당 단어의 대답을 입력해주세요.",
				Required:    true,
			},
		},
	},
	Aliases: []string{"공부"},
	DetailedDescription: &DetailedDescription{
		Usage: "머핀아 배워 (등록할 단어) (대답)",
		Examples: []string{"머핀아 배워 안녕 안녕!",
			"머핀아 배워 \"야 죽을래?\" \"아니요 ㅠㅠㅠ\"",
			"머핀아 배워 미간은_누구야? 이봇의_개발자요",
		},
	},
	Category: Chattings,
	MessageRun: func(ctx *MsgContext) {
		learnRun(ctx.Command, ctx.Session, ctx.Msg, &ctx.Args)
	},
	ChatInputRun: func(ctx *ChatInputContext) {
		var args *[]string
		learnRun(ctx.Command, ctx.Session, ctx.Inter, args)
	},
}

func addPrefix(arr []string) (newArr []string) {
	for _, item := range arr {
		newArr = append(newArr, fmt.Sprintf("- %s", item))
	}
	return
}

func learnRun(c *Command, s *discordgo.Session, m any, args *[]string) {
	var userId, command, result string

	igCommands := []string{}
	switch m := m.(type) {
	case *discordgo.MessageCreate:
		userId = m.Author.ID

		if len(*args) < 2 {
			s.ChannelMessageSendEmbedReply(m.ChannelID, &discordgo.MessageEmbed{
				Title:       "❌ 오류",
				Description: "올바르지 않ㅇ은 용법이에요.",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "사용법",
						Value:  utils.InlineCode(c.DetailedDescription.Usage),
						Inline: true,
					},
					{
						Name:  "예시",
						Value: utils.CodeBlockWithLanguage("md", strings.Join(addPrefix(c.DetailedDescription.Examples), "\n")),
					},
				},
				Color: int(utils.EFail),
			}, m.Reference())
			return
		}

		command = strings.ReplaceAll((*args)[0], "_", " ")
		result = strings.ReplaceAll((*args)[1], "_", " ")
	case *discordgo.InteractionCreate:
		s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})

		userId = m.Member.User.ID

		optsMap := map[string]*discordgo.ApplicationCommandInteractionDataOption{}
		for _, opt := range m.ApplicationCommandData().Options {
			optsMap[opt.Name] = opt
		}

		if opt, ok := optsMap["단어"]; ok {
			command = opt.StringValue()
		}

		if opt, ok := optsMap["대답"]; ok {
			result = opt.StringValue()
		}
	}

	for _, command := range Discommand.Commands {
		igCommands = append(igCommands, command.Name)
		igCommands = append(igCommands, command.Aliases...)
	}

	ignores := []string{"미간", "Migan", "migan", "간미"}
	ignores = append(ignores, igCommands...)

	disallows := []string{
		"@everyone",
		"@here",
		fmt.Sprintf("<@%s>", configs.Config.Bot.OwnerId),
	}

	for _, ig := range ignores {
		if strings.Contains(command, ig) {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ 오류",
				Description: "해ㄷ당 단어는 배우기 껄끄ㄹ럽네요.",
				Color:       int(utils.EFail),
			}

			switch m := m.(type) {
			case *discordgo.MessageCreate:
				s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
			case *discordgo.InteractionCreate:
				s.InteractionResponseEdit(m.Interaction, &discordgo.WebhookEdit{
					Embeds: &[]*discordgo.MessageEmbed{embed},
				})
			}
			return
		}
	}

	for _, di := range disallows {
		if strings.Contains(result, di) {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ 오류",
				Description: "해당 단ㅇ어의 대답으로 하기 좀 그렇ㄴ네요.",
				Color:       int(utils.EFail),
			}

			switch m := m.(type) {
			case *discordgo.MessageCreate:
				s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
			case *discordgo.InteractionCreate:
				s.InteractionResponseEdit(m.Interaction, &discordgo.WebhookEdit{
					Embeds: &[]*discordgo.MessageEmbed{embed},
				})
			}
		}
	}

	_, err := databases.Learns.InsertOne(context.TODO(), databases.InsertLearn{
		Command:   command,
		Result:    result,
		UserId:    userId,
		CreatedAt: time.Now(),
	})
	if err != nil {
		fmt.Println(err)
		embed := &discordgo.MessageEmbed{
			Title:       "❌ 오류",
			Description: "단어를 배우는데 오류가 생겼어요.",
			Color:       int(utils.EFail),
		}

		switch m := m.(type) {
		case *discordgo.MessageCreate:
			s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
		case *discordgo.InteractionCreate:
			s.InteractionResponseEdit(m.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{embed},
			})
		}
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "✅ 성공",
		Description: fmt.Sprintf("%s 배웠어요.", hangul.GetJosa(command, hangul.EUL_REUL)),
		Color:       int(utils.ESuccess),
	}

	switch m := m.(type) {
	case *discordgo.MessageCreate:
		s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
	case *discordgo.InteractionCreate:
		s.InteractionResponseEdit(m.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		})
	}
}
