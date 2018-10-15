# fragcenter
Stream display system to view active streams on an rtmp server.

Running fragcenter

### Get fragcenter locally

You can get fragcenter by running `go get github.com/fragforce/fragcenter`

You can run fragcenter by either running it with `go run fragcenter.go` in the src directory.  
You can also build it locally running `go build github.com/fragforce/fragcenter`

There are flags that can be set on startup.

`-host`  
    Set the host that is running the rtmp server (default is 127.0.0.1)  

`-port`  
    Set the port the rtmp server is serving on. (default is 8080)

`-web`  
    Set the port fragcenter uses to host it's own web server (default is 3000)  

This Also means the nginx.conf in the repo is copied to `/srv/rtmp/nginx.conf`. Please move it whereever you want and adjust the docker command accordingly.

### Docker command to run the rtmp server we built this for
This will default to port 1935 (default rtmp port) and port 8080 (default for stats pages)  
`docker run -it -d --rm -p 1935:1935 -p 8080:80 -v /srv/rtmp/nginx.conf:/opt/nginx/nginx.conf alfg/nginx-rtmp`

### Examples:  
1 host: Single host for rtmp server/fragcenter/web browser using the previous docker command   
    `fragcenter`

1 host: rtmp server/fragcenter custom web port  
    `fragcenter -host=<external_ip_of_host> -web=<port_to_host_web_pages_on>`

2 host: rtmp server customer stats port, fragcenter server custom web port  
    `fragcenter -host=<external_ip_of_host> -port=<stats_page_port> -web=<port_to_host_web_pages_on>`