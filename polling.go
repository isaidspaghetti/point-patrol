package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func StartPolling() {
	// Load environment variables
	//envErr := godotenv.Load()
	//rapidAPIKey := os.Getenv("RAPIDAPI_KEY")
	//rapidAPIHost := os.Getenv("RAPIDAPI_HOST")
	//if envErr != nil {
	//	log.Fatal("Error loading .env file")
	//}
	log.Println("test")

	ticker := time.NewTicker(6 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("tick")
			events, err := fetchEvents()
			if err != nil {
				log.Println("Error fetching events:", err)
				continue
			}
			BroadcastEvents(events)
		}
	}
}

func fetchEvents() ([]byte, error) {
	// create new request
	nbaURL := "https://allsportsapi2.p.rapidapi.com/api/basketball/matches/live"
	req, err := http.NewRequest("GET", nbaURL, nil)
	if err != nil {
		return nil, err
	}

	// Add required headers (manually)
	req.Header.Add("x-rapidapi-host", "allsportsapi2.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "3a7ee1f2bemsha1fc0b8d2bc760bp1d78e8jsn2d555b06efd5")
	// Add required headers (auto)
	//    rapidAPIHost := "allsportsapi2.p.rapidapi.com"
	//    rapidAPIKey := os.Getenv("RAPIDAPI_KEY")
	//    if rapidAPIKey == "" {
	//        log.Println("RAPIDAPI_KEY environment variable is not set")
	//        return nil
	//    }
	//    req.Header.Add("x-rapidapi-host", rapidAPIHost)
	//    req.Header.Add("x-rapidapi-key", rapidAPIKey)

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	return body, nil
	//var events EventResponse
	//if err := json.Unmarshal(body, &events); err != nil {
	//	return nil, err
	//}
	//fmt.Printf("body: %s\n", body)
	//
	//fmt.Printf("EVENTS", &events)
	//return &events, nil

}

//func broadcastEvents(events *EventResponse) {
//	// process each event individually
//	if events == nil {
//		log.Println("No events from external api")
//		return
//	}
//
//	for _, event := range events.Events {
//		eventData, err := json.Marshal(event)
//		if err != nil {
//			log.Println("Error serializing eventData:", err)
//			continue
//		}
//
//		// Extract team slugs from event
//		eventTeamSlugs := map[string]bool{
//			event.HomeTeam.Slug: true,
//			event.AwayTeam.Slug: true,
//		}
//		fmt.Printf("team slugs: %v\n", eventTeamSlugs)
//
//		// Broadcast this event to relevant clients
//		BroadcastEvent(eventData, eventTeamSlugs)
//	}
//}

//func broadcastEvents(events []byte) {
//	if len(events) == 0 {
//		log.Println("No events from external API")
//		return
//	}
//
//	// Broadcast the raw JSON event data to clients
//	clientsMutex.Lock()
//	for client := range clients {
//		err := client.Conn.WriteMessage(websocket.TextMessage, events)
//		if err != nil {
//			log.Printf("Error broadcasting to client: %v", err)
//		}
//	}
//	clientsMutex.Unlock()
//}
