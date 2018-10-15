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
`docker run -it -d --rm -p 1935:1935 -p 8080:80 -v /srv/rtmp/nginx.conf:/opt/nginx/nginx.conf alfg/nginx-rtmp`

### Examples:  
Single host for rtmp server/fragcenter/web browser using the previous docker command   
    `fragcenter`

If the host you are running the rtmp server is the one you are also running fragcenter but the web browser is on a different computer.  
    `fragcenter -host=<external_ip_of_host> -web=<port_to_host_web_pages_on>`

If the host you are running the rtmp server is the one you are also running fragcenter but the web browser is on a different computer.  
    `fragcenter -host=<external_ip_of_host> -web=<port_to_host_web_pages_on>`