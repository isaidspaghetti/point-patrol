package main

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Client struct {
	Conn       *websocket.Conn
	StationIDs map[string]bool
}

var (
	clients      = make(map[*Client]bool)
	clientsMutex sync.Mutex
)

// SubscribeClient adds client to the subscription list
func SubscribeClient(stationIDs []string, conn *websocket.Conn) {
	// Convert slice of ids to a map
	stationIDMap := make(map[string]bool)
	for _, id := range stationIDs {
		stationIDMap[id] = true
	}

	client := &Client{
		Conn:       conn,
		StationIDs: stationIDMap,
	}

	clientsMutex.Lock()
	clients[client] = true
	clientsMutex.Unlock()

	// handle client disconnection with a goroutine listening
	go func() {
		defer func() {
			clientsMutex.Lock()
			delete(clients, client)
			clientsMutex.Unlock()
			conn.Close()
		}()

		for {
			_, _, err := conn.NextReader()
			if err != nil {
				break
			}
		}
	}()

}

// IsTeamBroadcastedByStation checks if the station's team matches the team in the event
func IsTeamBroadcastedByStation(stationIDs map[string]bool, teamSlug string) bool {
	for stationID := range stationIDs {
		mappedTeamSlug, exists := stationToTeamMap[stationID]
		if !exists {
			continue
		}
		if mappedTeamSlug == teamSlug {
			return true
		}
	}
	return false
}

// BroadcastEvent sends event data to subscribed clients
func BroadcastEvent(eventData []byte, eventTeamSlugs map[string]bool) {
	clientsMutex.Lock()
	clientsCopy := make(map[*Client]bool)
	for client := range clients {
		clientsCopy[client] = true
	}
	clientsMutex.Unlock()

	var wg sync.WaitGroup
	for client := range clientsCopy {
		wg.Add(1)
		go func(c *Client) {
			defer wg.Done()
			for teamSlug := range eventTeamSlugs {
				if IsTeamBroadcastedByStation(c.StationIDs, teamSlug) {
					err := c.Conn.WriteMessage(websocket.TextMessage, eventData)
					if err != nil {
						log.Printf("Error writing to WebSocket: %v", err)
					}
					break
				}
			}
		}(client)
	}
	wg.Wait()
}
