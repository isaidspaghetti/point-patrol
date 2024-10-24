//package main
//
//import (
//	"github.com/tidwall/gjson"
//	"io"
//	"log"
//	"net/http"
//	"time"
//)
//
//var externalAPIEndpoints = []string{
//	"https://allsportsapi2.p.rapidapi.com/api/american-football/matches/live",
//	"https://allsportsapi2.p.rapidapi.com/api/basketball/matches/live",
//	"https://allsportsapi2.p.rapidapi.com/api/ice-hockey/matches/live",
//}
//
//// startDataFetching starts the central timer to fetch data from the external score API periodically
//func startDataFetching() {
//	ticker := time.NewTicker(9 * time.Second)
//	defer ticker.Stop()
//
//	for {
//		<-ticker.C
//		fetchAndProcessData()
//	}
//}
//
//// fenchAndProcessData fetches data from the api and processes it
//func fetchAndProcessData() {
//	apiResponse := getNBAScores() // start with just NBA scores
//	if apiResponse == nil {
//		return
//	}
//
//	processAPIResponse(apiResponse)
//}
//
//func getNBAScores() []byte {
//	// create new request
//	nbaURL := "https://allsportsapi2.p.rapidapi.com/api/basketball/matches/live"
//	req, err := http.NewRequest("GET", nbaURL, nil)
//	if err != nil {
//		return nil
//	}
//
//	// Add required headers (manually)
//	req.Header.Add("x-rapidapi-host", "allsportsapi2.p.rapidapi.com")
//	req.Header.Add("x-rapidapi-key", "3a7ee1f2bemsha1fc0b8d2bc760bp1d78e8jsn2d555b06efd5")
//	// Add required headers (auto)
//	//    rapidAPIHost := "allsportsapi2.p.rapidapi.com"
//	//    rapidAPIKey := os.Getenv("RAPIDAPI_KEY")
//	//    if rapidAPIKey == "" {
//	//        log.Println("RAPIDAPI_KEY environment variable is not set")
//	//        return nil
//	//    }
//	//    req.Header.Add("x-rapidapi-host", rapidAPIHost)
//	//    req.Header.Add("x-rapidapi-key", rapidAPIKey)
//
//	// Make request
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		log.Println("Error making API request:", err)
//		return nil
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		log.Printf("API Request failed with status code %d %s", resp.StatusCode, resp.Status)
//		return nil
//	}
//
//	// read the response body
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Println("Error reading response body:", err)
//		return nil
//	}
//
//	return body
//}
//
//// Process the api response and trigger message broadcasts
//func processAPIResponse(responseBody []byte) {
//	// Parse the events areay from teh JSON response
//	events := gjson.GetBytes(responseBody, "events")
//	if !events.Exists() {
//		log.Println("No events found in API response")
//		return
//	}
//
//	// Build a mapping of team slugs to event ids for active events
//	newTeamToEvent := make(map[string]int)
//	eventMessages := make(map[int][]byte) // Map of eventID to JSON data
//
//	// Iterate over each event in the vents array
//	events.forEach(func(key, value gjson.Result) bool {
//		// check if event is in progress
//	})
//}
