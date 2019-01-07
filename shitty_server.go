package main

import(
	"fmt"
	"os"
	"flag"
	"log"
	"path/filepath"

	"net"
	"net/http"
)

const(
	DEFAULT_HTTP_SERVER_PORT int = 8080

	HOSTNAME_ERROR_MESSAGE string = "HOSTNAME_ERROR"
	LOCALIP_ERROR_MESSAGE string = "LOCALIP_ERROR"
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
		log.Printf("getting hostname: %v\n", err)
		hostname = HOSTNAME_ERROR_MESSAGE
	}
	host_url := fmt.Sprintf("http://%s:%d", hostname, *port)

	localIP, err := getLocalIP()
	if err != nil {
		log.Printf("getting local IP: %v\n", err)

		localIP = LOCALIP_ERROR_MESSAGE
	}


	log.Printf("Server running on %s (%s)\n", host_url, localIP)

	log.Fatal(http.ListenAndServe( fmt.Sprintf(":%d", *port), nil ))
}

// Get preferred outbound ip of this machine.
// Code by Mr.Wang from Next Door,
// from https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}
