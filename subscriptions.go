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

// BroadcastEvent sends event data to subscribed clients. v7
func BroadcastEvent(eventData []byte, eventTeamSlugs map[string]bool) {
	clientsMutex.Lock()
	clientsCopy := make([]*Client, 0, len(clients))
	for client := range clients {
		clientsCopy = append(clientsCopy, client)
	}
	clientsMutex.Unlock()

	var wg sync.WaitGroup
	for _, client := range clientsCopy {
		wg.Add(1)
		go func(c *Client) {
			defer wg.Done()
			for teamSlug := range eventTeamSlugs {
				if IsTeamBroadcastedByStation(c.StationIDs, teamSlug) {
					err := c.Conn.WriteMessage(websocket.TextMessage, eventData)
					if err != nil {
						log.Printf("Error writing to client %v: %v", c.Conn.RemoteAddr(), err)
						// Remove the client on error.
						clientsMutex.Lock()
						delete(clients, c)
						clientsMutex.Unlock()
						c.Conn.Close()
					}
					break
				}
			}
		}(client)
	}
	wg.Wait()
}
