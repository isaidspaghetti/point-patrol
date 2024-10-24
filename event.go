package main

type EventResponse struct {
	Events []GameEvent `json:"events"`
}

type EventData struct {
	LiveMatches []GameEvent `json:"live_matches"`
}

// Provide some structs to get our hometeam and awayteam slugs
type GameEvent struct {
	MatchID   string `json:"match_id"`
	HomeTeam  Team   `json:"homeTeam"`
	AwayTeam  Team   `json:"awayTeam"`
	ID        int    `json:"id"`
	Score     string `json:"score"`
	EventTime string `json:"event_time"`
}

type Team struct {
	Slug string `json:"slug"`
}
