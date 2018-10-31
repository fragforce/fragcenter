package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

//LiveStreams is the datastructure of the stats xml page
type LiveStreams struct {
	Applications []struct {
		Name string `xml:"name"`
		Live []struct {
			Stream struct {
				Name string `xml:"name"`
				BWIn int    `xml:"bw_in"`
			} `xml:"stream"`
		} `xml:"live"`
	} `xml:"server>application"`
}

func main() {
	streamHost := flag.String("host", "127.0.0.1", "Host that the rtmp server is running on.")
	streamPort := flag.String("port", "8080", "Port the rtmp server is outputting http traffic")
	webPort := flag.String("web", "3000", "Port the webserver runs on.")
	pollInterval := flag.Int("poll", 10, "Polling interval")

	flag.Parse()

	fmt.Printf("Monitoring RTMP host %s:%s for live streams.\n", *streamHost, *streamPort)

	fmt.Printf("Starting stats checker, polling every %d seconds.\n", *pollInterval)
	go statsCheck(*streamHost, *streamPort, *pollInterval)

	fmt.Printf("Fragcenter is now running on port %s. Hit 'ctrl + c' to stop.\n", *webPort)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(fmt.Sprintf(":%s", *webPort), nil)
}

func marshalLiveStream(body []byte) (*LiveStreams, error) {
	var streams LiveStreams
	err := xml.Unmarshal(body, &streams)
	if err != nil {
		return nil, err
	}

	return &streams, nil
}

func statsCheck(host, port string, pollInterval int) {
	url := fmt.Sprintf("http://%s:%s/stats", host, port)
	for {
		fmt.Println("Checking Stats...")
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal("Problem getting stats page.\n", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic("Couldn't read the response body.")
		}

		liveStreams, err := marshalLiveStream(body)
		if err != nil {
			panic("Couldn't marshal the XML body to a struct.")
		}

		var active []string

		for _, application := range liveStreams.Applications {
			if application.Name == "stream" {
				for _, live := range application.Live {
					if live.Stream.BWIn == 0 {
						fmt.Printf("Stream '%s' is stopped. Ignoring.\n", live.Stream.Name)
						continue
					}
					active = append(active, live.Stream.Name)
					fmt.Printf("Found live stream '%s'.\n", live.Stream.Name)
				}
			}
		}

		sort.Strings(active)
		writeHTML(active, host, port)

		time.Sleep(time.Duration(pollInterval) * time.Second)
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

	for count, name := range streams {
		if count < 3 {
			bodyLines = append(bodyLines, strings.Replace(strings.Replace(strings.Replace(baseVideo, "<streamName>", name, -1), "<stereamHost>", host, -1), "<streamPort>", port, -1))
		}
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
