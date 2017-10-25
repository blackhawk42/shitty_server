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
		fmt.Fprintf(os.Stderr, "use: %s [-p port]", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	
	var port = flag.Int("p", DEFAULT_HTTP_SERVER_PORT, "`port` to be used for HTTP server")
	var root_dir = flag.String("d", ".", "`root directory` to serve from")
	
	flag.Parse()
	// Main
	
	
	http.Handle("/", http.FileServer(http.Dir(*root_dir)))
	log.Printf("Server running on port %d\n", *port)
	
	log.Fatal(http.ListenAndServe( fmt.Sprintf(":%d", *port), nil ))
}
