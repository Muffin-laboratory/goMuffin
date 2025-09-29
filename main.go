package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
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
		cmds = append(cmds, cmd.ApplicationCommand)
	}

	_, err = dg.ApplicationCommandBulkOverwrite(dg.State.User.ID, "", cmds)
	if err != nil {
		log.Println(err)
	}

	defer databases.GetDatabase().Disconnect()

	if port := &configs.GetConfig().IntegrateMDC.Server.Port; *port != 0 {
		go func() {
			log.Printf("[goMuffin] Muffin debug console 사용을 위한 서버가 포트 %d로 열렸어요.", *port)
			if err := server.Start(fmt.Sprintf(":%d", *port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalln(err)
			}
		}()
	} else {
		log.Println("[goMuffin] Muffin debug console 사용을 위한 서버가 꺼졌어요.")
	}

	defer server.Close()

	log.Println("[goMuffin] 봇이 실행되고 있어요. 버전:", configs.MuffinVersion)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
