//package main
//
//import (
//	"log"
//	"sync"
//)
//
//// The hub manages clients, subscriptions, and message broadcasts
//type Hub struct {
//	clients       map[*Client]bool            // Registered clients
//	register      chan *Client                // Channel for registering clients
//	unregister    chan *Client                // Channel for unregistering clients
//	subscriptions map[string]map[*Client]bool // Map of station IDs to clients subscribed to them
//	broadcast     chan *Message               // Channel for incoming messages to broadcast
//	mu            sync.RWMutex                // Mutex to protect shared resources
//}
//
//// Initialize the hub
//func NewHub() *Hub {
//	return &Hub{
//		clients:       make(map[*Client]bool),
//		register:      make(chan *Client),
//		unregister:    make(chan *Client),
//		subscriptions: make(map[string]map[*Client]bool),
//		broadcast:     make(chan *Message),
//	}
//}
//
//// global hub instance
//var hub = NewHub()
//
//// Main loop for the hub
//func (h *Hub) run() {
//	for {
//		select {
//		// Add a client
//		case client := <-h.register:
//			h.mu.Lock()
//			// Add client to hub
//			h.clients[client] = true
//			// Add client to stations (ids were added to client in the route handler)
//			for stationID := range client.stations {
//				if h.subscriptions[stationID] == nil {
//					h.subscriptions[stationID] = make(map[*Client]bool)
//				}
//				h.subscriptions[stationID][client] = true
//			}
//			h.mu.Unlock()
//			log.Printf("Client %s registered", client.clientID)
//
//		case client := <-h.unregister:
//			h.mu.Lock()
//			if _, ok := h.clients[client]; ok { // if that client is in the hub...
//				// Remove client from hub
//				delete(h.clients, client)
//				// Remove client from subscriptions
//				for stationID := range client.stations {
//					delete(h.subscriptions[stationID], client)
//					if len(h.subscriptions[stationID]) == 0 {
//						delete(h.subscriptions, stationID)
//					}
//				}
//				client.connection.Close()
//				log.Printf("Client %s unregistered", client.clientID)
//			}
//			h.mu.Unlock()
//
//		case message := <-h.broadcast:
//			// Broadcast a message to all clients subscribed to the station ID
//			h.mu.RLock()
//			clients := h.subscriptions[message.StationID]
//			for client := range clients {
//				client.writeMessage(message.Data)
//			}
//			h.mu.RUnlock()
//		}
//	}
//}
//
////
////func
////sendTestScoreUpdates()
////{
////	ticker := time.NewTicker(5 * time.Second)
////	defer ticker.Stop()
////
////	for {
////		select {
////		case t := <-ticker.C:
////			// simulate score updates fo rmultiple games
////			for gameID := 1; gameID <= 2; gameID++ { // assume 2 games for now
////				scoreUpdate := fmt.Sprintf("Game ID: %d, Score Update at %s", gameID, t.Format("15.04.05"))
////				message := &Message{
////					GameID: gameID,
////					Data:   []byte(scoreUpdate),
////				}
////				hub.broadcast <- message
////			}
////		},,,, 7u
////	}
////}
