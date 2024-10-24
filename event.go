package main

type EventResponse struct {
	Events []GameEvent `json:"events"`
}

// Provide some structs to get our hometeam and awayteam slugs
type GameEvent struct {
	HomeTeam Team `json:"homeTeam"`
	AwayTeam Team `json:"awayTeam"`
	ID       int  `json:"id"`
}

type Team struct {
	Slug string `json:"slug"`
}
