package scripts

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"git.wh64.net/muffin/goMuffin/configs"
	"github.com/bwmarrin/discordgo"
)

func DeleteAllCommands() {
	var answer string
	id := flag.String("id", "", "디스코드 봇의 토큰")

	flag.Parse()

	fmt.Printf("정말로 모든 명령어를 삭제하시겠어요? [y/N]: ")
	fmt.Scanf("%s", &answer)
	if strings.ToLower(answer) != "y" && strings.ToLower(answer) != "yes" {
		os.Exit(1)
	}

	if *id == "" {
		panic(fmt.Errorf("--id 플래그의 값이 필요해요."))
	}

	c := http.Client{}
	req, err := http.NewRequest("PUT", discordgo.EndpointApplicationGlobalCommands(*id), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bot "+configs.Config.Bot.Token)

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
