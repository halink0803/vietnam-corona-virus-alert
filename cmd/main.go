package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	crawler "github.com/halink0803/corona-alerts-bot/news-crawler"
	"github.com/halink0803/corona-alerts-bot/news-crawler/storage/sqlite"
	cli "github.com/urfave/cli/v2"
)

const (
	botKeyFlag    string = "bot-key"
	sleepDuration        = 1 * time.Minute
)

var (
	supportedChannels = []string{"689719240462172313"}
)

func main() {
	app := cli.NewApp()
	app.Name = "Corona Virus Alert bot"
	app.Action = run

	app.Flags = append(
		app.Flags,
		&cli.StringFlag{
			Name:    botKeyFlag,
			Usage:   "key for the bot",
			EnvVars: []string{"BOT_KEY"},
		},
	)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	botKey := c.String(botKeyFlag)
	sugar, flush, err := NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()
	crawler := crawler.NewCrawler(sugar)
	storage, err := sqlite.NewSqliteStorage()
	if err != nil {
		return err
	}
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

	latestNews, err := storage.GetLatestNews()
	if err != nil {
		return err
	}

	go func() {
		var (
			err  error
			news string
		)
		for {
			news, err = crawler.Start()
			if err != nil {
				log.Println(err)
				break
			}
			if latestNews != news {
				latestNews = news
				err := storage.StoreNews(news)
				if err != nil {
					log.Println("failed to store news", err)
					break
				}
				sendMessageToChannel(dg, latestNews)
				log.Println(latestNews)
			}
			time.Sleep(sleepDuration)
		}
	}()

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

	const (
		AboutMessage = `
		>>> **Corona Virus Alert About**
Send alert about new cases in Vietnam
Source from Vietnam Ministry of Health: https://ncov.moh.gov.vn
Bot invite link
<https://discordapp.com/api/oauth2/authorize?client_id=689005737015377920&permissions=18432&scope=bot>
		`
	)

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "cva!about" {
		log.Println(m.ChannelID)
		if _, err := s.ChannelMessageSend(m.ChannelID, AboutMessage); err != nil {
			log.Println(err)
		}
	}

	if m.Content == "cva!subscribe" {
		supportedChannels = append(supportedChannels, m.ChannelID)
		log.Printf("subscribed channel %s to alert", m.ChannelID)
	}
}

func sendMessageToChannel(s *discordgo.Session, message string) {
	// make it quote
	message = ">>> " + message
	for _, channelID := range supportedChannels {
		_, err := s.ChannelMessageSend(channelID, message)
		if err != nil {
			log.Println(err)
		}
	}
}
