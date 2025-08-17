package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"github.com/bwmarrin/discordgo"
)

func main() {
	err := dg.Open()
	if err != nil {
		log.Println("[goMuffin] 봇을 시작할 수 없어요.")
		log.Fatalln(err)
	}

	defer dg.Close()

	// 봇의 상태메세지 변경
	go func() {
		for {
			dg.UpdateCustomStatus("ㅅ살려주세요..!")
			time.Sleep(time.Minute * 10)
		}
	}()

	cmds := []*discordgo.ApplicationCommand{}
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

		cmds = append(cmds, cmd.ApplicationCommand)
	}

	_, err = dg.ApplicationCommandBulkOverwrite(dg.State.User.ID, "", cmds)
	if err != nil {
		log.Println(err)
	}

	defer databases.GetDatabase().Disconnect()

	log.Println("[goMuffin] 봇이 실행되고 있어요. 버전:", configs.MuffinVersion)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
