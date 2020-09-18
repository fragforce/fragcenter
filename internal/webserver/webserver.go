package webserver

import (
	"fmt"
	"log"
	"net/http"
)

func StartWebServer(webPort string) (err error) {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", webPort), nil); err != nil {
		log.Fatal(err)
	}
	return
}