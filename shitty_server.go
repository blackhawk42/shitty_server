package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"net"
	"net/http"
)

const (
	// Default port to bind the HTTP server to.
	DEFAULT_HTTP_SERVER_PORT int = 8080

	HOSTNAME_ERROR_MESSAGE string = "HOSTNAME_ERROR"
	LOCALIP_ERROR_MESSAGE  string = "LOCALIP_ERROR"
)

func main() {
	// Flag config
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "use: %s [-p port] [-d root_directory]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	var port = flag.Int("p", DEFAULT_HTTP_SERVER_PORT, "`port` to be used for HTTP server")
	var rootDir = flag.String("d", ".", "`root directory` to serve from")

	flag.Parse()

	// Main

	http.Handle("/", http.FileServer(http.Dir(*rootDir)))

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("getting hostname: %v\n", err)
		hostname = HOSTNAME_ERROR_MESSAGE
	}
	hostURL := fmt.Sprintf("http://%s:%d", hostname, *port)

	localIP, err := externalIP()
	if err != nil {
		log.Printf("getting local IP: %v\n", err)
	}
	hostIPURL := fmt.Sprintf("http://%s:%d", localIP, *port)

	log.Printf("Server running on %s (%s)\n", hostURL, hostIPURL)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

// Loop interfaces for a suitable IP address.
//
// Thank you kindly to Sebastian from:
// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return LOCALIP_ERROR_MESSAGE, err
	}

	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return LOCALIP_ERROR_MESSAGE, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			return ip.String(), nil
		}
	}

	return "127.0.0.1", fmt.Errorf("No suitable ip address found, possible lack of network connection; returning standard loopback address")
}
