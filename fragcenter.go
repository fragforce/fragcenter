package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type stream struct {
	name   string
	url    *url.URL
	mpdURL *url.URL
}

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
	streamHost := envOrFlagStr("STREAMHOST", "host", "127.0.0.1", "Host that the rtmp server is running on.")
	intStreamHost := envOrFlagStr("INTSTREAMHOST", "intHost", "127.0.0.1", "Internal container that the rtmp server is running on.")
	streamPort := envOrFlagStr("STREAMPORT", "port", "8080", "Port the rtmp server is outputting http traffic")
	webPort := envOrFlagStr("WEBPORT", "web", "3000", "Port the webserver runs on.")
	pollIntervalStr := envOrFlagStr("POLL", "poll", "10", "Polling interval")
	appName := envOrFlagStr("APPNAME", "appname", "stream", "Stream application name")
	streamPullKey := envOrFlagStr("STREAMPULLKEY", "pullKey", "bogus", "Stream key to use for auth'ing stream pulls")

	flag.Parse()

	pollInterval, err := strconv.Atoi(*pollIntervalStr)
	if err != nil {
		fmt.Println("Poll interval is not an integer. Using the default.")
		pollInterval = 10
	}

	// Wait for the RTMP server to come up
	time.Sleep(2 * time.Second)

	fmt.Printf("Monitoring RTMP host %s:%s for live streams.\n", *streamHost, *streamPort)

	fmt.Printf("Starting stats checker, polling every %d seconds for streams named '%s'.\n", pollInterval, *appName)
	go statsCheck(*streamHost, *intStreamHost, *streamPort, pollInterval, *appName, *streamPullKey)

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

func makeStreamURL(host string, port string, appName string, streamPullKey string, name string) *url.URL {
	u := url.URL{
		Scheme: "rtmp",
		Host:   host + ":" + port,
		Path:   fmt.Sprintf("%s/%s", appName, name),
	}

	q := u.Query()
	q.Set("key", streamPullKey)
	u.RawQuery = q.Encode()

	return &u
}

func statsCheck(host, intHost, port string, pollInterval int, appName string, streamPullKey string) {
	url := fmt.Sprintf("http://%s:%s/stats", intHost, port)
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
					fmt.Printf("Found live stream '%s'\n", stream.Name)
				}
			}
		}

		sort.Strings(active)

		streams := make([]stream, 0)

		for _, name := range active {
			r := stream{
				name: name,
				url:  makeStreamURL(host, port, appName, streamPullKey, name),
			}
			streams = append(streams, r)
		}

		if err := writeHTML(streams, host, port); err != nil {
			fmt.Printf("Problem running write html: %s\n", err)
		}

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

func writeHTML(streams []stream, host string, port string) error {

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
    <a href="<streamURL>"><video data-dashjs-player autoplay muted src="http://<streamHost>:<streamPort>/dash/<streamName>/index.mpd"></video></a>
    <br/>
    <q><streamName></q>
  </div>`

	for count, s := range streams {
		if count < 3 {
			bodyLines = append(bodyLines, strings.Replace(strings.Replace(strings.Replace(strings.Replace(baseVideo, "<streamName>", s.name, -1), "<streamHost>", host, -1), "<streamPort>", port, -1), "<streamURL>", s.url.String(), -1))
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
