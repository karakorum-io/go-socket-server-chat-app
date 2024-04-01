package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

var ip, port string

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	http.ServeFile(w, r, "home.html")
}

func main() {

	fmt.Println("Server IP: ")
	getInput(&ip)

	if !isValidIP(ip) {
		log.Fatal("Invalid IP!")
	}

	fmt.Println("Server PORT: ")
	getInput(&port)

	if !isValidPort(port) {
		log.Fatal("Invalid Port!")
	}

	log.Println("Starting server @ " + ip + ":" + port)

	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/public", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	server := &http.Server{
		Addr:              ip + ":" + port,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func getInput(value *string) {
	fmt.Scan(value)
}

func isValidIP(ipStr string) bool {
	return net.ParseIP(ipStr) != nil
}

func isValidPort(portStr string) bool {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false // not a number
	}
	return port > 0 && port <= 65535
}
