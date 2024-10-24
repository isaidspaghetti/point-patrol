package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func StartPolling() {
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

// fetchEvents fetches live events from the SportsDataAPI.
func fetchEvents() (*EventResponse, error) {
	apiURL := "https://allsportsapi2.p.rapidapi.com/api/basketball/matches/live"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request failed: %w", err)
	}

	// Retrieve API key from environment variable.
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("RAPIDAPI_KEY")
	if apiKey == "" {
		return nil, errors.New("RAPIDAPI_KEY environment variable is not set")
	}

	// Add required headers.
	req.Header.Add("x-rapidapi-host", "allsportsapi2.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", apiKey)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-200 status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body failed: %w", err)
	}

	var events EventResponse
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("unmarshalling JSON failed: %w", err)
	}

	log.Printf("Fetched %d events", len(events.Events))
	return &events, nil
}

// broadcastEvents processes and broadcasts each event to relevant clients.
func broadcastEvents(events *EventResponse) {
	if events == nil || len(events.Events) == 0 {
		log.Println("No events to broadcast")
		return
	}

	for _, rawEvent := range events.Events {
		// Extract team slugs by unmarshalling into TempGameEvent.
		var tempEvent GameEvent
		if err := json.Unmarshal(rawEvent, &tempEvent); err != nil {
			log.Printf("Error unmarshalling TempGameEvent: %v", err)
			continue
		}

		// Extract team slugs.
		eventTeamSlugs := map[string]bool{
			tempEvent.HomeTeam.Slug: true,
			tempEvent.AwayTeam.Slug: true,
		}

		// Broadcast the raw JSON to relevant clients.
		log.Printf("Broadcasting event ID %d to relevant clients", tempEvent.ID)
		BroadcastEvent(rawEvent, eventTeamSlugs)
	}
}
