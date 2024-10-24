//package main
//
//import (
//	"encoding/json"
//	"github.com/gorilla/websocket"
//	"log"
//	"sync"
//	"time"
//)
//
//// Constants for WebSocket read/write timeouts
//const (
//	pongWait   = 60 * time.Second
//	pingPeriod = (pongWait * 9) / 10
//	writeWait  = 10 * time.Second
//)
//
//// Client is a connected user
//type Client struct {
//	connection         *websocket.Conn
//	stationID          string          // radio station they are currently listening to
//	eventSubscriptions map[string]bool // eventIDs to subscribe to
//	clientID           string
//	mu                 sync.Mutex
//}
//
//// Message to be broadcasted to clients
//type Message struct {
//	StationID string // Station id to which the message pertains
//	Data      []byte
//}
//
//// subscribeToStation subscribes the client to a station ID
//func (c *Client) subscribeToStation(stationID string) {
//	// determine team names for station
//	teamSlug := stationToTeamSlug(stationID)
//	if teamSlug == "" {
//		log.Printf("Client %s tried to subscribe to station without a team associated: %s", c.clientID, stationID)
//		return
//	}
//	c.mu.Lock()
//
//	// Check for a live event for that team
//	eventID := teamSlugToEvent(teamSlug)
//	if eventID == "" {
//		logPrintf("No event is currently live for the team: %s", teamSlug)
//	}
//	// add the eventID to the client's subscriptions
//	if !c.eventSubscriptions[stationID] {
//		c.eventSubscriptions[stationID] = true
//		log.Printf("Client %s subscribed to event %s (team %s)", c.clientID, eventID, teamSlug)
//		// Register interest in this team with the EventManager
//		// update the hub howevdr you think we should
//	}
//	c.mu.Unlock()
//}
//
//// unsubscribeFromStation unsubscribes the client from a station ID
//func (c *Client) unsubscribeFromStation(stationID string) {
//	teamSlug := stationToTeamSlug(stationID)
//	if teamSlug == "" {
//		log.Printf("Client %s tried to unsubscribe from invalid station ID: %s", c.clientID, stationID)
//		return
//	}
//	c.mu.Lock()
//	if c.subscriptions[stationID] {
//		delete(c.subscriptions, stationID)
//		delete(c.teamSlugs, teamSlug)
//		log.Printf("Client %s unsubscribed from station %s", c.clientID, stationID)
//		// Remove interest in this team from the EventManager
//		eventManager.removeClientFromTeam(c, teamSlug)
//	}
//	c.mu.Unlock()
//}
//
//// writeMessage sends a message to the client
//func (c *Client) writeMessage(message []byte) {
//	c.connection.SetWriteDeadline(time.Now().Add(writeWait))
//	if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
//		log.Println("Error writing message to client %s: %v:", c.clientID, err)
//		c.connection.Close()
//	}
//}
//
//// startPingTimer sends periodic pings to keep connection alive
//func (c *Client) startPingTimer() {
//	ticker := time.NewTicker(pingPeriod)
//	defer func() {
//		ticker.Stop()
//		c.connection.Close()
//	}()
//	for {
//		select {
//		case <-ticker.C:
//			c.connection.SetWriteDeadline(time.Now().Add(writeWait))
//			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
//				log.Printf("Error sending ping to client %s: $v", c.clientID, err)
//				return
//			}
//		}
//	}
//}
//
//// readPump reads message from the WebSocket connection
//func (c *Client) readPump() {
//	defer func() {
//		hub.unregister <- c
//		c.connection.Close()
//	}()
//	c.connection.SetReadLimit(512)
//	c.connection.SetReadDeadline(time.Now().Add(pongWait))
//	c.connection.SetPongHandler(func(string) error {
//		c.connection.SetReadDeadline(time.Now().Add(pongWait))
//		return nil
//	})
//	for {
//		_, message, err := c.connection.ReadMessage()
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//				log.Printf("Client %s disconnected unexpectedly: %v", c.clientID, err)
//			}
//			break
//		}
//		// Handle subscription updates from the client
//		c.handleClientMessage(message)
//	}
//}
//
//// handleClientMessage processes messages received from the client (subscriptions)
//func (c *Client) handleClientMessage(message []byte) {
//	// assume the client can send a subscription message in JSON format
//	var msg struct {
//		Action    string `json:"action"`    // "subscribe" or "unsubscribe"
//		StationID string `json:"stationID"` // Station ID to subscribe/unsub to
//	}
//	if err := json.Unmarshal(message, &msg); err != nil {
//		log.Printf("Invalid message from client %s: %s", c.clientID, err)
//		return
//	}
//
//	switch msg.Action {
//	case "subscribe":
//		c.subscribeToStation(msg.StationID)
//	case "unsubscribe":
//		c.unsubscribeFromStation(msg.StationID)
//	default:
//		log.Printf("Unknown action from client %s: %s", c.clientID, msg.Action)
//	}
//}
//
//// subscribeTostation cubscribes the client to a station ID
//func (c *Client) subscribeToStation(stationID string) {
//	// Check team exists for station id
//	team, exists := stationToTeamMap[stationID]
//	if !exists {
//		log.Printf("Client %s tried to subscribe to invalid station ID: %s", c.clientID, stationID)
//		return
//	}
//
//	// Subscribe the client
//	hub.mu.Lock()
//	if !c.stations[stationID] { // do we actually need to check? just set it?
//		c.stations[stationID] = true
//		if hub.subscriptions[stationID] == nil {
//			hub.subscriptions[stationID] = make(map[*Client]bool)
//		}
//		hub.subscriptions[stationID][c] = true
//		log.Printf("Client %s subscribed to station %s", c.clientID, stationID)
//	}
//	hub.mu.Unlock()
//}
//
//// unsubscribeFromStation unsubs the client from the station ID
//func (c *Client) unsubscribeFromStation(stationID string) {
//	hub.mu.Lock()
//	if c.stations[stationID] {
//		delete(c.stations, stationID)
//		delete(hub.subscriptions[stationID], c)
//		if len(hub.subscriptions[stationID]) == 0 {
//			delete(hub.subscriptions, stationID)
//		}
//		log.Printf("Client %s unsubscribed from station %s", c.clientID, stationID)
//	}
//	hub.mu.Unlock()
//}
//
////
////func (c *Client) writePump() {
////	// Set a ping timer to keep connection alive
////	ticker := time.NewTicker(pingPeriod)
////	defer func() {
////		ticker.Stop()
////		c.connection.Close()
////	}()
////
////	for {
////		select {
////		case message, ok := <-c.send:
////			c.connection.SetWriteDeadline(time.Now().Add(writeWait))
////			if !ok {
////				// the hub closed the channel.
////				c.connection.WriteMessage(websocket.CloseMessage, []byte{})
////				return
////			}
////
////			// Drain the queue to get the latest message
////			for {
////				select {
////				case latestMessage := <-c.send:
////					message = latestMessage // Update message to latest
////				default:
////					// No more messages in the queue, exit draining
////					break
////				}
////			}
////
////			// Create a new writer for the WebSocket Message
////			w, err := c.connection.NextWriter(websocket.TextMessage)
////			if err != nil {
////				log.Printf("Error getting websocket writer: %v", err)
////				return
////			}
////
////			w.Write(message)
////
////			// Add queued messages to the current websocket message
////			n := len(c.send)
////			for i := 0; i < n; i++ {
////				w.Write([]byte("\n"))
////				w.Write(<-c.send)
////			}
////
////			// Close the writer
////			if err := w.Close(); err != nil {
////				log.Printf("Error closing websocket writer: %v", err)
////				return
////			}
////		case <-ticker.C:
////			// Send a ping message to keep the connection alive
////			c.connection.SetWriteDeadline(time.Now().Add(writeWait))
////			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
////				log.Printf("Error sending ping: %v", err)
////				return
////			}
////		}
////	}
////}
