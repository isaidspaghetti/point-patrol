package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Implement proper origin checking in production
	},
}

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameter for stationIDs
	stationIDsParam := r.URL.Query().Get("stationIDs")
	if stationIDsParam == "" {
		http.Error(w, "stationID is required", http.StatusBadRequest)
		return
	}

	// Split the stationIDs into a slice
	stationIDs := strings.Split(stationIDsParam, ",")

	// Upgrade the HTTP connection to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	// Subscribe the client to updates based on stationIDs
	SubscribeClient(stationIDs, conn)
}
