package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"flag"

	obsws "github.com/christopher-dG/go-obs-websocket"
)

var ()

//LiveStreams is the datastructure of the stats xml page
type LiveStreams struct {
	Applications []struct {
		Name string `xml:"name"`
		Live struct {
			Streams []struct {
				Name string `xml:"name"`
				BWIn int    `xml:"bw_in"`
			} `xml:"stream"`
		} `xml:"live"`
	} `xml:"server>application"`
}

// use environment variables to set the default values of flags
func envOrFlagStr(envName, flagName, flagDefault, usage string) *string {
	if value, exists := os.LookupEnv(envName); exists {
		return flag.String(flagName, value, usage)
	} else {
		return flag.String(flagName, flagDefault, usage)
	}
}

func main() {
	// stream host flags
	streamHost := envOrFlagStr("STREAMHOST", "host", "127.0.0.1", "Host that the rtmp server is running on.")
	streamPort := envOrFlagStr("STREAMPORT", "port", "80", "Port the rtmp server is outputting http traffic")
	appName := envOrFlagStr("APPNAME", "appname", "stream", "Stream application name")
	primaryStream := envOrFlagStr("PRIMARYSTREAM", "primary", "superstream", "Stream key to check for")
	// obs host flags
	obsHost := envOrFlagStr("OBSHOST", "obshost", "127.0.0.1", "Host that the obs application is running on.")
	obsPort := envOrFlagStr("OBSPORT", "obsport", "4444", "Port the obs websocket server listening on.")
	primaryScene := envOrFlagStr("PRIMARYSCENE", "primaryscene", "Scene", "Primary Scene to show while stream is up")
	backupScene := envOrFlagStr("BACKUPSCENE", "backupscene", "Backup", "Scene to switch to if main stream is down")

	flag.Parse()

	var currentScene string

	// loop to check if primary stream is active
	for {
		log.Printf("Checking Stats for %s", *primaryStream)
		active, err := statsCheck(*streamHost, *streamPort, *appName, *primaryStream)
		if err != nil {
			log.Fatal(err)
		}

		if !active && (currentScene == *primaryScene || currentScene == "") {
			log.Printf("primaryscene is down switching to backup")
			if err := switchToScene(*obsHost, *obsPort, *backupScene); err != nil {
				log.Fatal(err)
			}
			currentScene = *backupScene
		} else if active && (currentScene == *backupScene || currentScene == "") {
			log.Printf("primaryscene is up switching to primary")
			if err := switchToScene(*obsHost, *obsPort, *primaryScene); err != nil {
				log.Fatal(err)
			}
			currentScene = *primaryScene
		} else if active {
			log.Printf("primary already active no switching")
		} else {
			log.Printf("backup already active no switching")
		}

		time.Sleep(time.Second * 5)
	}

}

func statsCheck(intHost, port, appName, primary string) (active bool, err error) {
	url := fmt.Sprintf("http://%s:%s/stats", intHost, port)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	current, err := marshalLiveStream(body)
	if err != nil {
		return false, err
	}

	for _, app := range current.Applications {
		if app.Name == appName {
			for _, stream := range app.Live.Streams {
				if stream.Name == primary {
					if stream.BWIn == 0 {
						return false, nil
					} else {
						return true, nil
					}
				}
			}
		}
	}

	return false, err
}

func switchToScene(host, port, scene string) (err error) {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return
	}

	// Connect a client.
	obsClient := obsws.Client{Host: host, Port: portInt}
	if err := obsClient.Connect(); err != nil {
		return err
	}

	defer obsClient.Disconnect()

	setSceneReq := obsws.NewSetCurrentSceneRequest(scene)
	if setSceneResp, err := setSceneReq.SendReceive(obsClient); err != nil {
		return err
	} else {
		fmt.Println(setSceneResp.Status_)
	}

	obsClient.Disconnect()

	return
}

func marshalLiveStream(body []byte) (*LiveStreams, error) {
	var streams LiveStreams
	err := xml.Unmarshal(body, &streams)
	if err != nil {
		return nil, err
	}

	return &streams, nil
}
