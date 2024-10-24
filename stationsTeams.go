package main

type TeamStations struct {
	Team     string   `json:"team"`
	Stations []string `json:"stations"`
}

// NBA teams
var stationToTeamMap = map[string]string{
	"wqaq":  "miami-heat",
	"kgmz":  "golden-state-warriors",
	"wmfs":  "memphis-grizzlies",
	"wscr":  "chicago-bulls",
	"wxyt":  "detroit-pistons",
	"wwj":   "detroit-pistons",
	"wfan":  "brooklyn-nets",
	"wfan2": "brooklyn-nets",
	"wzgc":  "atlanta-hawks",
	"wtem":  "washington-wizards",
	"wwl":   "new-orleans-pelicans",
	"test":  "botafogo",
}

// stationsToTeamSlug retursn t5he team slug for a given station id
func stationsToTeamSlug(stationID string) string {
	return stationToTeamMap[stationID]
}
