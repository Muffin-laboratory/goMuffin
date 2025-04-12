package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	var answer string
	id := flag.String("id", "", "discordBot's token")

	flag.Parse()

	fmt.Printf("Do you want to delete all commands? [y/N]: ")
	fmt.Scanf("%s", &answer)
	if strings.ToLower(answer) != "y" {
		os.Exit(1)
	}

	if os.Getenv("BOT_TOKEN") == "" {
		panic(fmt.Errorf("You need a BOT_TOKEN environment."))
	}

	if *id == "" {
		panic(fmt.Errorf("You need a --id flag value."))
	}

	c := http.Client{}
	req, err := http.NewRequest("PUT", discordgo.EndpointApplicationGlobalCommands(*id), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bot "+os.Getenv("BOT_TOKEN"))

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
