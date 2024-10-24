package main

import (
	"encoding/json"
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
			broadcastEvents(events)
		}
	}
}

func fetchEvents() (*EventResponse, error) {
	// create new request
	nbaURL := "https://allsportsapi2.p.rapidapi.com/api/basketball/matches/live"
	req, err := http.NewRequest("GET", nbaURL, nil)
	if err != nil {
		return nil, err
	}

	// Add required headers (manually)
	req.Header.Add("x-rapidapi-host", "allsportsapi2.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "hidden")
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
		return nil, err
	}

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("received non-200 status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var events EventResponse
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, err
	}
	fmt.Printf("body: %s\n", body)

	fmt.Printf("EVENTS", &events)
	return &events, nil
}

func broadcastEvents(events *EventResponse) {
	// process each event individually
	if events == nil {
		log.Println("No events from external api")
		return
	}

	for _, event := range events.Events {
		eventData, err := json.Marshal(event)
		if err != nil {
			log.Println("Error serializing eventData:", err)
			continue
		}

		// Extract team slugs from event
		eventTeamSlugs := map[string]bool{
			event.HomeTeam.Slug: true,
			event.AwayTeam.Slug: true,
		}
		fmt.Printf("team slugs: %v\n", eventTeamSlugs)

		// Broadcast this event to relevant clients
		BroadcastEvent(eventData, eventTeamSlugs)
	}
}
