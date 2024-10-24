//package main
//
//import (
//	"sync"
//)
//
//// Event represents data for a single event
//type Event struct {
//	ID       int
//	HomeTeam string
//	AwayTeam string
//	Status   string
//}
//
//type EventManager struct {
//	events      map[int]*Event // Event ID to Event Data
//	teamToEvent map[string]int // match team slug to current event id
//	mu          sync.RWMutex
//}
//
//func NewEventManager() *EventManager {
//	return &EventManager{
//		events:      make(map[int]*Event),
//		teamToEvent: make(map[string]int),
//	}
//}
//
//var eventManager = NewEventManager()
//
//// GetEventIDForTeam returns current event ID for a team or 0 if not playing
//func (em *EventManager) GetEventIDForTeam(teamSlug string) int {
//	em.mu.RLock()
//	defer em.mu.RUnlock()
//	return em.teamToEvent[teamSlug]
//}
//
//func (em *EventManager) AddEvent(event *Event) {
//	// add event to the vent manager
//	em.events[event.ID] = event
//	// add the team mappings of the new event
//	em.teamToEvent[event.HomeTeam] = event.ID
//	em.teamToEvent[event.AwayTeam] = event.ID
//}
