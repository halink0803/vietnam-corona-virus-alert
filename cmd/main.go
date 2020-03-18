package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	crawler "github.com/halink0803/corona-alerts-bot/news-crawler"
	cli "github.com/urfave/cli/v2"
)

const (
	botKeyFlag = "bot-key"
)

func main() {
	app := cli.NewApp()
	app.Name = "Corona Virus Alert bot"
	app.Action = run

	app.Flags = append(
		app.Flags,
		&cli.StringFlag{
			Name:    "bot-key",
			Usage:   "key for the bot",
			EnvVars: []string{"BOT_KEY"},
		},
	)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	botKey := c.String("bot-key")
	fmt.Println(botKey)
	sugar, flush, err := NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()
	crawler := crawler.NewCrawler(sugar)
	crawler.Start()
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + botKey)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return err
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return err
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	return dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
