package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/commands/dev"
	"git.wh64.net/muffin/goMuffin/components"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/handler"
	"git.wh64.net/muffin/goMuffin/modals"
	"git.wh64.net/muffin/goMuffin/scripts"
	"github.com/bwmarrin/discordgo"
	"github.com/devproje/commando"
	"github.com/devproje/commando/types"
)

func init() {
	go commands.GetDiscommand().LoadCommand(commands.HelpCommand)
	go commands.GetDiscommand().LoadCommand(commands.DataLengthCommand)
	go commands.GetDiscommand().LoadCommand(commands.LearnCommand)
	go commands.GetDiscommand().LoadCommand(commands.LearnedDataListCommand)
	go commands.GetDiscommand().LoadCommand(commands.InformationCommand)
	go commands.GetDiscommand().LoadCommand(commands.DeleteLearnedDataCommand)
	go commands.GetDiscommand().LoadCommand(dev.ReloadPromptCommand)
	go commands.GetDiscommand().LoadCommand(commands.SwitchModeCommand)
	go commands.GetDiscommand().LoadCommand(commands.ChatCommand)
	go commands.GetDiscommand().LoadCommand(commands.RegisterCommand)
	go commands.GetDiscommand().LoadCommand(commands.DeregisterCommand)
	go commands.GetDiscommand().LoadCommand(dev.BlockCommand)
	go commands.GetDiscommand().LoadCommand(dev.UnblockCommand)

	go commands.GetDiscommand().LoadComponent(components.DeleteLearnedDataComponent)
	go commands.GetDiscommand().LoadComponent(components.PaginationEmbedComponent)
	go commands.GetDiscommand().LoadComponent(components.RegisterComponent)
	go commands.GetDiscommand().LoadComponent(components.DeregisterComponent)
	go commands.GetDiscommand().LoadComponent(components.SelectChatComponent)
	go commands.GetDiscommand().LoadComponent(components.DeleteChatComponent)

	go commands.GetDiscommand().LoadModal(modals.PaginationEmbedModal)
}

func main() {
	command := commando.NewCommando(os.Args[1:])
	config := configs.GetConfig()

	if len(os.Args) > 1 {
		command.Root("delete-all-commands", "봇의 모든 슬래시 커맨드를 삭제합니다.", scripts.DeleteAllCommands,
			types.OptionData{
				Name: "id",
				Desc: "봇의 디스코드 아이디",
				Type: types.STRING,
			},
			types.OptionData{
				Name:  "isYes",
				Short: []string{"y"},
				Type:  types.BOOLEAN,
			},
		)

		err := command.Execute()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	dg, _ := discordgo.New("Bot " + config.Bot.Token)

	go dg.AddHandler(handler.MessageCreate)
	go dg.AddHandler(handler.InteractionCreate)

	err := dg.Open()
	if err != nil {
		log.Println("[goMuffin] 봇을 시작할 수 없어요.")
		log.Fatalln(err)
	}

	chatbot.New(dg)

	defer dg.Close()

	// 봇의 상태메세지 변경
	go func() {
		for {
			dg.UpdateCustomStatus("ㅅ살려주세요..!")
			time.Sleep(time.Minute * 10)
		}
	}()

	for _, cmd := range commands.GetDiscommand().Commands {
		if cmd.Name == commands.HelpCommand.Name {
			// 극한의 성능 똥망 코드 탄생!
			// 무려 똑같은 걸 반복해서 돌리는!
			for _, a := range commands.GetDiscommand().Commands {
				cmd.Options[0].Choices = append(cmd.Options[0].Choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  a.Name,
					Value: a.Name,
				})
			}
		}

		if !cmd.RegisterApplicationCommand {
			continue
		}

		go dg.ApplicationCommandCreate(dg.State.User.ID, "", cmd.ApplicationCommand)
	}

	defer databases.Disconnect()

	log.Println("[goMuffin] 봇이 실행되고 있어요. 버전:", configs.MUFFIN_VERSION)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
