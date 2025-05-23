package commands

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/LoperLee/golang-hangul-toolkit/hangul"
	"github.com/bwmarrin/discordgo"
)

var learnArguments = utils.InlineCode("{user.name}") + "\n" +
	utils.InlineCode("{user.mention}") + "\n" +
	utils.InlineCode("{user.globalName}") + "\n" +
	utils.InlineCode("{user.id}") + "\n" +
	utils.InlineCode("{user.createdAt}") + "\n" +
	utils.InlineCode("{user.joinedAt}") + "\n" +
	utils.InlineCode("{muffin.version}") + "\n" +
	utils.InlineCode("{muffin.updatedAt}") + "\n" +
	utils.InlineCode("{muffin.statedAt}") + "\n" +
	utils.InlineCode("{muffin.name}") + "\n" +
	utils.InlineCode("{muffin.id}")

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
		Usage: fmt.Sprintf("%s배워 (등록할 단어) (대답)", configs.Config.Bot.Prefix),
		Examples: []string{
			fmt.Sprintf("%s배워 안녕 안녕!", configs.Config.Bot.Prefix),
			fmt.Sprintf("%s배워 \"야 죽을래?\" \"아니요 ㅠㅠㅠ\"", configs.Config.Bot.Prefix),
			fmt.Sprintf("%s배워 미간은_누구야? 이봇의_개발자요", configs.Config.Bot.Prefix),
			fmt.Sprintf("%s배워 \"나의 아이디를 알려줘\" \"너의 아이디는 {user.id}야.\"", configs.Config.Bot.Prefix),
		},
	},
	Category: Chatting,
	MessageRun: func(ctx *MsgContext) {
		if len(*ctx.Args) < 2 {
			utils.NewMessageSender(ctx.Msg).
				AddEmbeds(&discordgo.MessageEmbed{
					Title:       "❌ 오류",
					Description: "올바르지 않ㅇ은 용법이에요.",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "사용법",
							Value:  utils.InlineCode(ctx.Command.DetailedDescription.Usage),
							Inline: true,
						},
						{
							Name:   "사용 가능한 인자",
							Value:  learnArguments,
							Inline: true,
						},
						{
							Name:  "예시",
							Value: utils.CodeBlock("md", strings.Join(addPrefix(ctx.Command.DetailedDescription.Examples), "\n")),
						},
					},
					Color: utils.EmbedFail,
				}).
				SetReply(true).
				Send()
			return
		}

		learnRun(ctx.Msg, ctx.Msg.Author.ID, strings.ReplaceAll((*ctx.Args)[0], "_", " "), strings.ReplaceAll((*ctx.Args)[1], "_", " "))
	},
	ChatInputRun: func(ctx *ChatInputContext) {
		ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		})

		var command, result string

		if opt, ok := ctx.Inter.Options["단어"]; ok {
			command = opt.StringValue()
		}

		if opt, ok := ctx.Inter.Options["대답"]; ok {
			result = opt.StringValue()
		}

		learnRun(ctx.Inter, ctx.Inter.Member.User.ID, command, result)
	},
}

func addPrefix(arr []string) (newArr []string) {
	for _, item := range arr {
		newArr = append(newArr, fmt.Sprintf("- %s", item))
	}
	return
}

func learnRun(m any, userId, command, result string) {
	igCommands := []string{}

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
			utils.NewMessageSender(m).
				AddEmbeds(
					&discordgo.MessageEmbed{
						Title:       "❌ 오류",
						Description: "해ㄷ당 단어는 배우기 껄끄ㄹ럽네요.",
						Color:       utils.EmbedFail,
					}).
				SetReply(true).
				Send()
			return
		}
	}

	for _, di := range disallows {
		if strings.Contains(result, di) {
			utils.NewMessageSender(m).
				AddEmbeds(
					&discordgo.MessageEmbed{
						Title:       "❌ 오류",
						Description: "해당 단ㅇ어의 대답으로 하기 좀 그렇ㄴ네요.",
						Color:       utils.EmbedFail,
					}).
				SetReply(true).
				Send()
			return
		}
	}

	_, err := databases.Database.Learns.InsertOne(context.TODO(), databases.InsertLearn{
		Command:   command,
		Result:    result,
		UserId:    userId,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println(err)

		utils.NewMessageSender(m).
			AddEmbeds(&discordgo.MessageEmbed{
				Title:       "❌ 오류",
				Description: "단어를 배우는데 오류가 생겼어요.",
				Color:       utils.EmbedFail,
			}).
			SetReply(true).
			Send()
		return
	}

	utils.NewMessageSender(m).
		AddEmbeds(&discordgo.MessageEmbed{
			Title:       "✅ 성공",
			Description: fmt.Sprintf("%s 배웠어요.", hangul.GetJosa(command, hangul.EUL_REUL)),
			Color:       utils.EmbedSuccess,
		}).
		SetReply(true).
		Send()
}
