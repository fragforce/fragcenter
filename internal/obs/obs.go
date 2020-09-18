package obs

import (
	"fmt"
	"log"

	obs "github.com/christopher-dG/go-obs-websocket"
)

var (
	obsClient obs.Client

	stopClient    = make(chan string)
	clientStarted = make(chan string)
)

func StartOBSClient() {
	//go StatsCheck("stor22.gpmidi.net", "8081", 10, "stream")

	// Connect a client.
	obsClient = obs.Client{Host: "localhost", Port: 4444}
	if err := obsClient.Connect(); err != nil {
		log.Fatal(err)
	}
	defer obsClient.Disconnect()

	clientStarted <- ""

	<-stopClient

	obsClient.Disconnect()
}

func StopOBSClient() {
	stopClient <- ""
}

func ChangeSceneTo(sceneName string) (err error) {
	setSceneReq := obs.NewSetCurrentSceneRequest(sceneName)
	if setSceneResp, err := setSceneReq.SendReceive(obsClient); err != nil {
		return err
	} else {
		fmt.Println(setSceneResp.Status_)
	}

	return
}

func GetCurrentSceneStream() (sceneName, streamURL string, err error) {
	getSceneReq := obs.NewGetCurrentSceneRequest()

	if currentScene, err := getSceneReq.SendReceive(obsClient); err != nil {
		return "", "", err
	} else {
		sceneSettingReq := obs.NewGetSourceSettingsRequest(currentScene.Sources[0].Name, currentScene.Sources[0].Type)
		if sceneSettingResp, err := sceneSettingReq.SendReceive(obsClient); err != nil {
			log.Fatal(err)
		} else {
			if sceneSettingResp.SourceSettings["input"] == nil {
				return currentScene.Name, "", err
			} else {
				streamURL = fmt.Sprint(sceneSettingResp.SourceSettings["input"])
			}
		}
	}

	return
}

func CreateSource(request obs.DuplicateSceneItemRequest) (err error) {
	newSourceReq := obs.NewDuplicateSceneItemRequest(
		request.FromScene,
		request.ToScene,
		request.Item,
		request.ItemName,
		request.ItemID,
	)
	newSourceReq.SendReceive(obsClient)

	return
}

func UpdateSourceSettings(request obs.SetSourceSettingsRequest) (err error) {
	newSourceReq := obs.NewSetSourceSettingsRequest(request.SourceName, request.SourceType, request.SourceSettings)

	if newSourceResp, err := newSourceReq.SendReceive(obsClient); err != nil {
		return err
	} else {
		fmt.Printf(newSourceResp.Status())
	}
	return
}
