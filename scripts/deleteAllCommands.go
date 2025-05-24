package scripts

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"git.wh64.net/muffin/goMuffin/configs"
	"github.com/bwmarrin/discordgo"
	"github.com/devproje/commando"
	"github.com/devproje/commando/option"
)

func DeleteAllCommands(n *commando.Node) error {
	var answer string

	id, err := option.ParseString(*n.MustGetOpt("id"), n)
	if err != nil {
		return err
	}

	yes, _ := option.ParseBool(*n.MustGetOpt("isYes"), n)
	if !yes {
		fmt.Printf("정말로 모든 명령어를 삭제하시겠어요? [y/N]: ")
		fmt.Scanf("%s", &answer)
		if strings.ToLower(answer) != "y" && strings.ToLower(answer) != "yes" {
			fmt.Println("모든 명령어 삭제를 취소했어요.")
			os.Exit(1)
		}
	}

	c := http.Client{}
	req, err := http.NewRequest("PUT", discordgo.EndpointApplicationGlobalCommands(id), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bot "+configs.Config.Bot.Token)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}
