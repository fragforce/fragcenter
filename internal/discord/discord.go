package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fragforce/fragcenter/internal/config"
)

var (
	stopDiscordServer    = make(chan string)
	discordServerStopped = make(chan string)

	discordLoad = make(chan string)

	dConf discordConfig
)

type discordConfig struct {
	Token string `json:"token"`
	Game  string `json:"game,omitempty"`
}

func Init(confDir string) (err error) {
	if !strings.HasSuffix(confDir, "/") {
		confDir = confDir + "/"
	}

	if err := config.LoadConfig(confDir+"discord.json", dConf); err != nil {
		return err
	}

	return
}

// This function will be called (due to AddHandler) when the bot receives
// the "ready" event from Discord.
func readyDiscord(dg *discordgo.Session, event *discordgo.Ready, game string) {
	// if there is an error setting the game log and return
	if err := dg.UpdateStatus(0, game); err != nil {
		return
	}
}

// This function will be called (due to AddHandler) every time a new
// message is created on any channel that the authenticated bot has access to.
func discordMessageHandler(dg *discordgo.Session, m *discordgo.MessageCreate) {

}

// when a shutdown is sent close out services properly
func StopDiscordBot() {
	log.Println("stopping discord connections")
	stopDiscordServer <- ""

	<-discordServerStopped

	log.Println("discord connections stopped")
}

func StartDiscordBot() {
	// Initializing Discord connection
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + dConf.Token)
	if err != nil {
		return
	}

	dg.AddHandler(readyDiscord)
	dg.AddHandler(discordMessageHandler)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
		return
	}

	bot, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	log.Print("Invite the bot to your server with https://discordapp.com/oauth2/authorize?client_id=" + bot.ID + "&scope=bot")

	discordLoad <- ""

	<-stopDiscordServer

	// properly send a shutdown to the discord server so the bot goes offline.
	if err := dg.Close(); err != nil {
		log.Println(err)
	}

	// return the shutdown signal
	discordServerStopped <- ""
}
