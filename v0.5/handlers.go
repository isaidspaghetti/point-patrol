//package main
//
//import (
//	"github.com/gorilla/websocket"
//	"log"
//	"net/http"
//	"strings"
//)
//
//// handles incoming websocket connections
//func handleConnections(w http.ResponseWriter, r *http.Request) {
//	log.Println("new client connection")
//
//	// upgrade the initial GET request to a WebSocket
//	ws, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Println("Websocket upgrade error:", err)
//		return
//	}
//
//	// Create new client for the request
//	client := &Client{
//		connection: ws,
//		stations:   make(map[string]bool),
//		clientID:   r.RemoteAddr, // use remote addr as client ID for now
//	}
//
//	// Extract station IDs from query parameters (comma-separated list)
//	stationIDsParam := r.URL.Query().Get("stationIDs")
//	if stationIDsParam != "" {
//		stationIDs := strings.Split(stationIDsParam, ",")
//		for _, stationID := range stationIDs {
//			stationID = strings.TrimSpace(stationID)
//			// check if a team matches for that stationID
//			team, exists := stationToTeamMap[stationID]
//			if exists {
//				client.stations[stationID] = true
//				log.Printf("Client %s subscribed to station %s (team %s)", client.clientID, stationID, team)
//			} else {
//				log.Printf("No team found for station %s", stationID)
//			}
//		}
//	} else {
//		log.Printf("Client %s did not provide any station IDs", client.clientID)
//	}
//
//	// Register the client with the hub
//	hub.register <- client
//
//	// Start the client's read and ping pumps
//	go client.startPingTimer() // ping to keep conn alive
//	go client.readPump()       // read for disconnects
//
//	// Ensure client is unregisterd on exit
//	//defer func() {
//	//	hub.unregister <- client
//	//	ws.Close()
//	//}()
//
//	//// Handle multiple game IDs
//	//gameIDsStr := r.URL.Query().Get("game")
//	//if gameIDsStr != "" {
//	//	gameIDs := strings.Split(gameIDsStr, ",")
//	//	for _, idStr := range gameIDs {
//	//		idStr = strings.TrimSpace(idStr)
//	//		if idStr == "" {
//	//			continue
//	//		}
//	//		gameID, err := strconv.Atoi(idStr)
//	//		if err == nil {
//	//			client.games[gameID] = true
//	//			log.Printf("Client %s subscribed to game ID: %d\n", client.clientID, gameID)
//	//		}
//	//	}
//	//}
//
//}
//
//var upgrader = websocket.Upgrader{
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//
//	// Allow origins
//	CheckOrigin: func(r *http.Request) bool {
//		return true
//		//return r.Host == "localhost"
//	},
//}
