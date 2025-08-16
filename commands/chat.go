package commands

import (
	subcommand "git.wh64.net/muffin/goMuffin/commands/subcommands/chat"
	"git.wh64.net/muffin/goMuffin/configs"
	"github.com/bwmarrin/discordgo"
)

type chatCommandType string

var (
	chatCommandChatting chatCommandType = "하기"
	chatCommandList     chatCommandType = "목록"
	chatCommandCreate   chatCommandType = "생성"
	chatCommandDelete   chatCommandType = "삭제"
)

var ChatCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "대화",
		Description: "이 봇이랑 대화해요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        string(chatCommandList),
				Description: "채팅 목록을 나열해요.",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        string(chatCommandCreate),
				Description: "새로운 채팅을 생성해요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "이름",
						Description: "채팅의 이름을 정해요. (25자 이내)",
						MaxLength:   25,
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        string(chatCommandChatting),
				Description: "이 봇이랑 대화해요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "내용",
						Description: "대화할 내용",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        string(chatCommandDelete),
				Description: "채팅을 삭제해요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "제목",
						Description: "채팅의 제목",
						Required:    true,
					},
				},
			},
		},
	},
	DetailedDescription: DetailedDescription{
		Usage: configs.AddPrefix("%s대화 (목록/생성) [채팅 이름]"),
		Examples: []string{
			configs.AddPrefix("%s대화 목록"),
			configs.AddPrefix("%s대화 생성 머핀 냠냠"),
		},
	},
	Category: Chatting,
	Flags:    CommandFlagsIsRegistered | CommandFlagsIsBlocked,
	Run: func(ctx *ChatInputContext) error {
		var cType chatCommandType
		var str string

		if opt, ok := ctx.Inter.Options[string(chatCommandChatting)]; ok {
			ctx.Inter.DeferReply(nil)
			cType = chatCommandChatting
			str = opt.Options[0].StringValue()
		} else if opt, ok := ctx.Inter.Options[string(chatCommandCreate)]; ok {
			ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			})
			cType = chatCommandCreate
			str = opt.Options[0].StringValue()
		} else if _, ok := ctx.Inter.Options[string(chatCommandList)]; ok {
			ctx.Inter.DeferReply(nil)
			cType = chatCommandList
		} else if _, ok := ctx.Inter.Options[string(chatCommandDelete)]; ok {
			ctx.Inter.DeferReply(nil)
			cType = chatCommandDelete
			str = opt.Options[0].StringValue()
		}
		return chatCommandRun(cType, ctx.Inter, ctx.Inter.User, str)
	},
}

func chatCommandRun(cType chatCommandType, m any, user *discordgo.User, contentOrName string) error {
	switch cType {
	case chatCommandChatting:
		return subcommand.Chat(m, user, contentOrName)
	case chatCommandCreate:
		return subcommand.Create(m, user, contentOrName)
	case chatCommandList:
		return subcommand.List(m, user)
	case chatCommandDelete:
		return subcommand.Delete(m, user, contentOrName)
	}
	return nil
}
