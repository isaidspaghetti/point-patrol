package main

import (
	"log"
	"net/http"
)

func main() {
	// Start the polling
	go StartPolling()

	// Set up HTTP server and route
	// Fake frontend for testing
	http.Handle("/", http.FileServer(http.Dir("./frontend-dupe")))

	// HTTP Websocket handler
	http.HandleFunc("/ws", WebSocketHandler)

	addr := ":8080"
	log.Printf("Server started on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}

}
