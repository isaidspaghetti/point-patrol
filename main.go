package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	// start the hub
	go hub.run()
	go sendScoreUpdates()

	// run a fake frontend
	http.Handle("/", http.FileServer(http.Dir("./frontend-dupe")))
	log.Println("Server started on :8080")

	// HTTP Websocket handler
	http.HandleFunc("/ws", handleConnections)

	// Run server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")

	// upgrade the http connection to a websocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Create new client for the request
	client := &Client{
		connection: ws,
		send:       make(chan []byte, 256),
		games:      make(map[int]bool),
		clientID:   r.RemoteAddr, // use remote addr as client ID for now
	}

	// Register the client
	hub.register <- client

	// Ensure client is unregisterd on exit
	defer func() {
		hub.unregister <- client
		ws.Close()
	}()

	// Handle multiple game IDs
	gameIDsStr := r.URL.Query().Get("game")
	if gameIDsStr != "" {
		gameIDs := strings.Split(gameIDsStr, ",")
		for _, idStr := range gameIDs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			gameID, err := strconv.Atoi(idStr)
			if err == nil {
				client.games[gameID] = true
				log.Printf("Client %s subscribed to game ID: %d\n", client.clientID, gameID)
			}
		}
	}

	// Start goroutines to read and write messages
	go client.writePump()
	client.readPump() // to detect disconnects
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
		// return r.Host == "localhost"
	},
}
