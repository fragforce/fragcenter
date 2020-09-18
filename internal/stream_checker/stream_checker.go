package stream_checker

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

//LiveStreams is the data structure of the stats xml page
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

func marshalLiveStream(body []byte) (*LiveStreams, error) {
	var streams LiveStreams
	err := xml.Unmarshal(body, &streams)
	if err != nil {
		return nil, err
	}

	return &streams, nil
}

// StatsCheck runs a check against the nginx rtmp server server xml file for stream stats.
func StatsCheck(host, port string, pollInterval int, appName string) {
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
			if application.Name == appName {
				for _, stream := range application.Live.Streams {
					if stream.BWIn == 0 {
						fmt.Printf("Stream '%s' is stopped. Ignoring.\n", stream.Name)
						continue
					}
					active = append(active, stream.Name)
					fmt.Printf("Found live stream '%s'.\n", stream.Name)
				}
			}
		}

		sort.Strings(active)
		if err := writeHTML(active, host, port); err != nil {
			log.Println(err)
		}

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func GetActiveStreams(host, port string, appName string) (streams []string, err error) {
	url := fmt.Sprintf("http://%s:%s/stats", host, port)
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

	for _, application := range liveStreams.Applications {
		if application.Name == appName {
			for _, stream := range application.Live.Streams {
				if stream.BWIn == 0 {
					fmt.Printf("Stream '%s' is stopped. Ignoring.\n", stream.Name)
					continue
				}
				streams = append(streams, stream.Name)
				fmt.Printf("Found live stream '%s'.\n", stream.Name)
			}
		}
	}
	return
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
    <a href="rtmp://<streamHost>/stream/<streamName>"><video data-dashjs-player autoplay muted src="http://<streamHost>:<streamPort>/dash/<streamName>/index.mpd"></video></a>
    <br/>
    <q><streamName></q>
  </div>`

	for count, name := range streams {
		if count < 3 {
			bodyLines = append(bodyLines, strings.Replace(strings.Replace(strings.Replace(baseVideo, "<streamName>", name, -1), "<streamHost>", host, -1), "<streamPort>", port, -1))
		}
	}

	htmlBody = strings.Join(bodyLines, "\n")

	fo, err := os.Create("./public/index.html")
	if err != nil {
		return err
	}
	defer fo.Close()

	if _, err := fo.WriteString(htmlStart + htmlBody + htmlEnd); err != nil {
		log.Println(err)
	}

	if err := fo.Close(); err != nil {
		log.Println(err)
	}

	return nil
}
