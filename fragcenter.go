package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	xj "github.com/basgys/goxml2json"
)

var ()

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

	streamHost := flag.String("host", "127.0.0.1", "Host that the rtmp server is running on.")
	streamPort := flag.String("port", "8080", "Port the rtmp server is outputting http traffic")
	webPort := flag.String("web", "3000", "Port the webserver runs on.")

	flag.Parse()

	fmt.Println("rtmp host: " + *streamHost + ":" + *streamPort)

	fmt.Println("Starting web host on port " + *webPort)
	go webHost(*webPort)
	fmt.Println("Starting stats checker")
	go statsCheck(*streamHost, *streamPort)

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

func webHost(port string) {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":"+port, nil)
}

func statsCheck(host string, port string) {
	for {
		fmt.Println("Checking Stats")
		resp, err := http.Get("http://" + host + ":" + port + "/stats")
		if err != nil {
			log.Fatal("Problem getting stats page.\n", err)
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
			if application.Name == "live" {
				for _, live := range application.Live.Stream {
					if live.BwIn == "0" {
						fmt.Println("stream is stopped")
						continue
					}
					active = append(active, live.Name)
					fmt.Println(live.Name)
				}
				sort.Strings(active)
				writeHTML(active, host, port)
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

func writeHTML(streams []string, host string, port string) error {

	var bodyLines []string

	htmlBody := ""

	htmlStart := `<html>
<head>
  <title>Fragcenter</title>
  <script src="https://cdn.dashjs.org/latest/dash.all.min.js"></script>
  <script src="http://ajax.googleapis.com/ajax/libs/jquery/1.11.1/jquery.min.js"></script>
  <script type="text/javascript">
    function getPage(){
      var result;
      $.ajax({
        url: 'index.html',
        type: 'get',
        async: false,
        success: function(data) {
          result = data;
        }
      });
      return result;
    }

    current = getPage();

    function checkChanges(){
      check = getPage();
      if ( check != current) {
        location.reload();
      };
    }
    setInterval(checkChanges, 10000);
  </script>
  <style>
	video {
	width: 100%;
	padding-left: 1%;
	padding-right: 1%;
	}
	#container {
	width: 30%;
	padding-left: 1%;
	padding-right: 1%;
	float: left;
	}
  </style>
</head>
<body style="background-color:slategray;">
<div align="center">
`

	htmlEnd := `
</div>
</body>
</html>`

	baseVideo := `  <div id="container">
    <video data-dashjs-player autoplay muted src="http://<stereamHost>:<streamPort>/dash/<streamName>/index.mpd"></video>
    <br/>
    <q><streamName></q>
  </div>`

	for _, name := range streams {
		bodyLines = append(bodyLines, strings.Replace(strings.Replace(strings.Replace(baseVideo, "<streamName>", name, -1), "<stereamHost>", host, -1), "<streamPort>", port, -1))
	}

	htmlBody = strings.Join(bodyLines, "\n")

	fo, err := os.Create("./public/index.html")
	if err != nil {
		return err
	}
	defer fo.Close()

	fo.WriteString(htmlStart + htmlBody + htmlEnd)

	fo.Close()

	return nil
}
