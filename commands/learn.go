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

var learnArguments = "> " + utils.InlineCode("{user.name}") + "\n" +
	"> " + utils.InlineCode("{user.mention}") + "\n" +
	"> " + utils.InlineCode("{user.globalName}") + "\n" +
	"> " + utils.InlineCode("{user.id}") + "\n" +
	"> " + utils.InlineCode("{user.createdAt}") + "\n" +
	"> " + utils.InlineCode("{user.joinedAt}") + "\n" +
	"> " + utils.InlineCode("{muffin.version}") + "\n" +
	"> " + utils.InlineCode("{muffin.updatedAt}") + "\n" +
	"> " + utils.InlineCode("{muffin.statedAt}") + "\n" +
	"> " + utils.InlineCode("{muffin.name}") + "\n" +
	"> " + utils.InlineCode("{muffin.id}")

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
		Usage: configs.AddPrefix("%s배워 (등록할 단어) (대답)"),
		Examples: []string{
			configs.AddPrefix("%s배워 안녕 안녕!"),
			configs.AddPrefix("%s배워 \"야 죽을래?\" \"아니요 ㅠㅠㅠ\""),
			configs.AddPrefix("%s배워 미간은_누구야? 이봇의_개발자요"),
			configs.AddPrefix("%s배워 \"나의 아이디를 알려줘\" \"너의 아이디는 {user.id}야.\""),
		},
	},
	Category:                   Chatting,
	RegisterApplicationCommand: true,
	RegisterMessageCommand:     true,
	Flags:                      CommandFlagsIsRegistered | CommandFlagsIsBlocked,
	MessageRun: func(ctx *MsgContext) error {
		if len(*ctx.Args) < 2 {
			utils.NewMessageSender(ctx.Msg).
				AddComponents(utils.GetErrorContainer(
					discordgo.TextDisplay{
						Content: "올바르지 않ㅇ은 용법이에요.",
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **사용법**\n> %s", ctx.Command.DetailedDescription.Usage),
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **예시**\n%s", strings.Join(utils.AddPrefix("> ", ctx.Command.DetailedDescription.Examples), "\n")),
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **사용 가능한 인자**\n%s", learnArguments),
					},
				)).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return nil
		}

		return learnRun(ctx.Msg, ctx.Msg.Author.ID, strings.ReplaceAll((*ctx.Args)[0], "_", " "), strings.ReplaceAll((*ctx.Args)[1], "_", " "))
	},
	ChatInputRun: func(ctx *ChatInputContext) error {
		err := ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			return err
		}

		var command, result string

		if opt, ok := ctx.Inter.Options["단어"]; ok {
			command = opt.StringValue()
		}

		if opt, ok := ctx.Inter.Options["대답"]; ok {
			result = opt.StringValue()
		}

		return learnRun(ctx.Inter, ctx.Inter.Member.User.ID, command, result)
	},
}

func learnRun(m any, userId, command, result string) error {
	igCommands := []string{}

	for _, command := range instance.Commands {
		igCommands = append(igCommands, command.Name)
		igCommands = append(igCommands, command.Aliases...)
	}

	ignores := []string{"미간", "Migan", "migan", "간미"}
	ignores = append(ignores, igCommands...)

	disallows := []string{
		"@everyone",
		"@here",
		fmt.Sprintf("<@%s>", configs.GetConfig().Bot.OwnerId),
	}

	for _, ig := range ignores {
		if strings.Contains(command, ig) {
			utils.NewMessageSender(m).
				AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "해ㄷ당 단어는 배우기 껄끄ㄹ럽네요."})).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return nil
		}
	}

	for _, di := range disallows {
		if strings.Contains(result, di) {
			utils.NewMessageSender(m).
				AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "해당 단ㅇ어의 대답으로 하기 좀 그렇ㄴ네요."})).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return nil
		}
	}

	if len([]rune(command)) > 100 {
		utils.NewMessageSender(m).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "단어는 100글자를 못 넘ㅇ어가요."})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
		return nil
	}

	_, err := databases.GetDatabase().Learns.InsertOne(context.TODO(), databases.Learn{
		Command:   command,
		Result:    result,
		UserId:    userId,
		CreatedAt: time.Now(),
	})
	if err != nil {
		utils.NewMessageSender(m).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "단어를 배우는데 오류가 생겼어요."})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
		return err
	}

	return utils.NewMessageSender(m).
		AddComponents(utils.GetSuccessContainer(
			discordgo.TextDisplay{
				Content: fmt.Sprintf("%s 배웠어요.", hangul.GetJosa(command, hangul.EUL_REUL)),
			},
		)).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
