package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/components"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/handler"
	"git.wh64.net/muffin/goMuffin/modals"
	"github.com/bwmarrin/discordgo"
)

func init() {
	go commands.Discommand.LoadCommand(commands.HelpCommand)
	go commands.Discommand.LoadCommand(commands.DataLengthCommand)
	go commands.Discommand.LoadCommand(commands.LearnCommand)
	go commands.Discommand.LoadCommand(commands.LearnedDataListCommand)
	go commands.Discommand.LoadCommand(commands.InformationCommand)
	go commands.Discommand.LoadCommand(commands.DeleteLearnedDataCommand)

	go commands.Discommand.LoadComponent(components.DeleteLearnedDataComponent)
	go commands.Discommand.LoadComponent(components.PaginationEmbedComponent)

	go commands.Discommand.LoadModal(modals.PaginationEmbedModal)
}

func main() {
	// command := commando.NewCommando(os.Args[1:])
	config := configs.Config

	// if len(os.Args) > 1 {
	// 	command.Root("delete-all-commands", "봇의 모든 슬래시 커맨드를 삭제합니다.", scripts.DeleteAllCommands,
	// 		types.OptionData{
	// 			Name: "id",
	// 			Desc: "봇의 디스코드 아이디",
	// 			Type: types.STRING,
	// 		},
	// 		types.OptionData{
	// 			Name:  "isYes",
	// 			Short: []string{"y"},
	// 			Type:  types.BOOLEAN,
	// 		},
	// 	)

	// 	command.Root("export", "머핀봇의 데이터를 추출합니다.", scripts.ExportData,
	// 		types.OptionData{
	// 			Name: "type",
	// 			Desc: "파일형식을 지정합니다. (json, jsonl, finetune)",
	// 			Type: types.STRING,
	// 		},
	// 		types.OptionData{
	// 			Name: "export-path",
	// 			Desc: "데이터를 저장할 위치를 지정합니다.",
	// 			Type: types.STRING,
	// 		},
	// 		types.OptionData{
	// 			Name: "refined",
	// 			Desc: "머핀 데이터를 있는 그대로 추출할 지, 가려내서 추출할 지를 지정합니다.",
	// 			Type: types.BOOLEAN,
	// 		},
	// 	)

	// 	err := command.Execute()
	// 	if err != nil {
	// 		_, _ = fmt.Fprintln(os.Stderr, err)
	// 		os.Exit(1)
	// 	}
	// 	return
	// }

	isAI := flag.Bool("ai", false, "시작할 때 AI모드로 설정합니다.")

	flag.Parse()

	dg, _ := discordgo.New("Bot " + config.Bot.Token)

	go dg.AddHandler(handler.MessageCreate)
	go dg.AddHandler(handler.InteractionCreate)

	err := dg.Open()
	if err != nil {
		log.Println("[goMuffin] 봇을 시작할 수 없어요.")
		log.Fatalln(err)
	}

	chatbot.New(dg)

	if *isAI {
		chatbot.ChatBot.SetMode(chatbot.ChatbotAI)
	}

	defer dg.Close()

	// 봇의 상태메세지 변경
	go func() {
		for {
			dg.UpdateCustomStatus("ㅅ살려주세요..!")
			time.Sleep(time.Minute * 10)
		}
	}()

	for _, cmd := range commands.Discommand.Commands {
		if cmd.Name == commands.HelpCommand.Name {
			// 극한의 성능 똥망 코드 탄생!
			// 무려 똑같은 걸 반복해서 돌리는!
			for _, a := range commands.Discommand.Commands {
				cmd.Options[0].Choices = append(cmd.Options[0].Choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  a.Name,
					Value: a.Name,
				})
			}
		}

		go dg.ApplicationCommandCreate(dg.State.User.ID, "", cmd.ApplicationCommand)
	}

	defer databases.Database.Client.Disconnect(context.TODO())

	log.Println("[goMuffin] 봇이 실행되고 있어요. 버전:", configs.MUFFIN_VERSION)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
