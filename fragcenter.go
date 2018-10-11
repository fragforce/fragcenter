package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	xj "github.com/basgys/goxml2json"
)

var (
	streamHost string
)

//LiveStreams are the active streams on the stats page.
type LiveStreams struct {
	Rtmp struct {
		NginxVersion     string `json:"nginx_version"`
		NginxRtmpVersion string `json:"nginx_rtmp_version"`
		Naccepted        string `json:"naccepted"`
		BytesIn          string `json:"bytes_in"`
		BwOut            string `json:"bw_out"`
		BytesOut         string `json:"bytes_out"`
		Compiler         string `json:"compiler"`
		Built            string `json:"built"`
		Pid              string `json:"pid"`
		Uptime           string `json:"uptime"`
		BwIn             string `json:"bw_in"`
		Server           struct {
			Application []struct {
				Name string `json:"name"`
				Live struct {
					Stream []struct {
						Time    string `json:"time"`
						BwIn    string `json:"bw_in"`
						BwOut   string `json:"bw_out"`
						BwVideo string `json:"bw_video"`
						Client  []struct {
							Flashver   string `json:"flashver"`
							Dropped    string `json:"dropped"`
							Avsync     string `json:"avsync"`
							Timestamp  string `json:"timestamp"`
							Active     string `json:"active"`
							ID         string `json:"id"`
							Address    string `json:"address"`
							Time       string `json:"time"`
							Swfurl     string `json:"swfurl,omitempty"`
							Publishing string `json:"publishing,omitempty"`
						} `json:"client"`
						Meta struct {
							Video struct {
								Height    string `json:"height"`
								FrameRate string `json:"frame_rate"`
								Codec     string `json:"codec"`
								Profile   string `json:"profile"`
								Compat    string `json:"compat"`
								Level     string `json:"level"`
								Width     string `json:"width"`
							} `json:"video"`
							Audio struct {
								Channels   string `json:"channels"`
								SampleRate string `json:"sample_rate"`
								Codec      string `json:"codec"`
								Profile    string `json:"profile"`
							} `json:"audio"`
						} `json:"meta"`
						Nclients   string `json:"nclients"`
						Publishing string `json:"publishing"`
						Name       string `json:"name"`
						BytesIn    string `json:"bytes_in"`
						BytesOut   string `json:"bytes_out"`
						BwAudio    string `json:"bw_audio"`
						Active     string `json:"active"`
					} `json:"stream"`
					Nclients string `json:"nclients"`
				} `json:"live"`
			} `json:"application"`
		} `json:"server"`
	} `json:"rtmp"`
}

func main() {
	go webHost()
	go statsCheck()

	fmt.Println("Fragcenter is now running. Send 'shutdown' or 'ctrl + c' to stop Fragcenter.")

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("cannot read from stdin")
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if line == "shutdown" {
			fmt.Println("Shutting down fragcenter.")
			return
		}
	}

}

func webHost() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":3000", nil)
}

func statsCheck() {
	for {
		resp, err := http.Get("http://" + streamHost + ":8080/stats")
		if err != nil {
			log.Fatal("", err)
		}

		body, err := ioutil.ReadAll(resp.Body)

		converted, err := xj.Convert(strings.NewReader(string(body)))
		if err != nil {
			panic("That's embarrassing...")
		}

		var active []string

		streams := LiveStreams{}
		json.Unmarshal(converted.Bytes(), &streams)

		for _, application := range streams.Rtmp.Server.Application {
			if application.Name == "stream" {
				fmt.Println("New batch")
				for _, live := range application.Live.Stream {
					if live.BwIn == "0" {
						continue
					}
					fmt.Println(live.Name)
					active = append(active, live.Name)
					fmt.Println(live.BwIn)
				}
				fmt.Println("")
				writeHTML(active)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func fileCheck() {
	for {
		time.Sleep(10 * time.Second)
		files, err := ioutil.ReadDir("/tmp/rtmp/active")
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			if !f.IsDir() {
				fmt.Println(strings.TrimSuffix(f.Name(), ".m3u8"))
			}
		}
	}
}

func writeHTML(streams []string) error {

	htmlBody := ``

	htmlStart := `<!DOCTYPE html>
<html>
<head>
    <title>Fragcenter</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/dashjs/2.4.1/dash.all.min.js"></script>
    <script src="https://cdn.dashjs.org/latest/dash.all.min.js"></script>
<style>
  	video {
	   width: 30%;
  	}
</style>
</head>
<body style="background-color:slategray;">
`

	htmlEnd := `
	<script type="text/javascript">
    function checkChanges(){
      $.get("index.html", function(data){
        if ( data != document.documentElement.innerHTML) {
          location.reload();
        };
      });
    }
    setInterval(checkChanges, 5000);
  </script>
</body>
</html>`

	baseVideo := `<video data-dashjs-player autoplay src="http://<stereamHost>:8080/dash/<streamname>/index.mpd" controls></video>`

	for count, name := range streams {
		if count < 3 {
			fmt.Println("There are less than 3 streams")
			fmt.Println(strings.Replace("<video data-dashjs-player autoplay src=\"http://<stereamHost>:8080/dash/<streamName>/index.mpd\" controls></video>", "<streamName>", name, -1))
		}
		//htmlBody = strings.Replace(BaseVideo, "<streamName>", name, -1)
	}

	fo, err := os.Create("./public/index.html")
	if err != nil {
		return err
	}
	defer fo.Close()

	writer := bufio.NewWriter(fo)

	fmt.Fprint(writer, htmlStart+htmlBody+htmlEnd)

	return nil
}
