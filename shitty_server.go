package main

import(
	"fmt"
	"os"
	"flag"
	"log"
	"path/filepath"
	"net/http"
)

const(
	DEFAULT_HTTP_SERVER_PORT int = 8080
)

func main() {
	// Flag config
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "use: %s [-p port] [-d root_directory]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	
	var port = flag.Int("p", DEFAULT_HTTP_SERVER_PORT, "`port` to be used for HTTP server")
	var root_dir = flag.String("d", ".", "`root directory` to serve from")
	
	flag.Parse()
	
	// Main
	
	http.Handle("/", http.FileServer(http.Dir(*root_dir)))
	
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("getting hostname: %v\n", err)
	}
	host_url := fmt.Sprintf("http://%s:%d", hostname, *port)
	log.Printf("Server running on %s\n", host_url)
	
	log.Fatal(http.ListenAndServe( fmt.Sprintf(":%d", *port), nil ))
}
