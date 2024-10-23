package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// The hub manages clients and their subscriptions
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client // Channel for registering new clients
	unregister chan *Client
	broadcast  chan *Message // Channel for broadcasting messages
	mu         sync.RWMutex  // A mutex to protect shared resources
}

// Initialize the hub
var hub = Hub{
	clients:    make(map[*Client]bool),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	broadcast:  make(chan *Message),
}

// Hub runner
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Println("Client registered")
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("Client unregistered")
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				if client.games[message.GameID] {
					select {
					case client.send <- message.Data:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

type Message struct {
	GameID int
	Data   []byte
}

func sendScoreUpdates() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			// simulate score updates fo rmultiple games
			for gameID := 1; gameID <= 2; gameID++ { // assume 2 games for now
				scoreUpdate := fmt.Sprintf("Game ID: %d, Score Update at %s", gameID, t.Format("15.04.05"))
				message := &Message{
					GameID: gameID,
					Data:   []byte(scoreUpdate),
				}
				hub.broadcast <- message
			}
		}
	}
}
