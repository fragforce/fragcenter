package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/fragforce/fragcenter/internal/discord"
	"github.com/fragforce/fragcenter/internal/obs"
	"github.com/fragforce/fragcenter/internal/stream_checker"
	"github.com/fragforce/fragcenter/internal/webserver"
)

var (
	shutdown    = make(chan string)
	servStopped = make(chan string)
)

// use environment variables to set the default values of flags
func envOrFlagStr(envName, flagName, flagDefault, usage string) *string {
	if value, exists := os.LookupEnv(envName); exists {
		return flag.String(flagName, value, usage)
	} else {
		return flag.String(flagName, flagDefault, usage)
	}
}

func main() {
	streamHost := envOrFlagStr("STREAMHOST", "host", "127.0.0.1", "Host that the rtmp server is running on.")
	streamPort := envOrFlagStr("STREAMPORT", "port", "8080", "Port the rtmp server is outputting http traffic")
	webPort := envOrFlagStr("WEBPORT", "web", "3000", "Port the web server runs on.")
	pollIntervalStr := envOrFlagStr("POLL", "poll", "10", "Polling interval")
	appName := envOrFlagStr("APPNAME", "appname", "stream", "Stream application name")
	confDir := envOrFlagStr("APPNAME", "appname", "stream", "Stream application name")
	flag.Parse()

	pollInterval, err := strconv.Atoi(*pollIntervalStr)
	if err != nil {
		fmt.Println("Poll interval is not an integer. Using the default.")
		pollInterval = 10
	}

	// Wait for the RTMP server to come up
	time.Sleep(2 * time.Second)

	fmt.Printf("starting obs websocket client")
	obs.StartOBSClient()

	fmt.Printf("Monitoring RTMP host %s:%s for live streams.\n", *streamHost, *streamPort)
	fmt.Printf("Starting stats checker, polling every %d seconds for streams named '%s'.\n", pollInterval, *appName)
	go stream_checker.StatsCheck(*streamHost, *streamPort, pollInterval, *appName)

	fmt.Printf("starting webserve to host fragcenter preview page\n")
	go webserver.StartWebServer(*webPort)

	fmt.Println("Starting discord session.")
	discord.Init(*confDir)
	go discord.StartDiscordBot()

	go catchSig()
	go console()

	fmt.Printf("Fragcenter is now running on port %s. Hit 'ctrl + c' to stop.\n", *webPort)

	<-shutdown

	discord.StopDiscordBot()

	fmt.Printf("Fragcenter has closed out\n")
}

func console() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("cannot read from stdin: %v\n", err)
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if line == "shutdown" {
			log.Println("shutting down the bot")
			log.Println("All services stopped")
			shutdown <- ""
			return
		}
	}
}

func catchSig() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		os.Interrupt)
	<-sigc
	log.Println("interrupt caught")

	shutdown <- ""
}
