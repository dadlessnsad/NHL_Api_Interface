package helpers

import (
	"encoding/json"
	"fmt"
	"sync"
)

type AllTeamsResponse struct {
	Teams []Team `json:"teams"`
}

type Team struct {
	Name   string         `json:"name"`
	Roster RosterResponse `json:"roster"`
}

type RosterResponse struct {
	Forwards   []Player `json:"forwards"`
	Defensemen []Player `json:"defensemen"`
	Goalies    []Player `json:"goalies"`
}

type Player struct {
	ID                  int    `json:"id"`
	Headshot            string `json:"headshot"`
	FirstName           Name   `json:"firstName"`
	FullName            string `json:"fullName"`
	LastName            Name   `json:"lastName"`
	SweaterNumber       int    `json:"sweaterNumber"`
	PositionCode        string `json:"positionCode"`
	ShootsCatches       string `json:"shootsCatches"`
	HeightInInches      int    `json:"heightInInches"`
	WeightInPounds      int    `json:"weightInPounds"`
	HeightInCentimeters int    `json:"heightInCentimeters"`
	WeightInKilograms   int    `json:"weightInKilograms"`
	BirthDate           string `json:"birthDate"`
	BirthCity           Name   `json:"birthCity"`
	BirthCountry        string `json:"birthCountry"`
}

type Name struct {
	Default string `json:"default"`
}

func FormatRosterData(rosterData string) (*RosterResponse, error) {
	var roster RosterResponse
	err := json.Unmarshal([]byte(rosterData), &roster)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling roster data: %v", err)
	}

	var wg sync.WaitGroup
	addFullName := func(p *Player) {
		defer wg.Done()
		p.FullName = p.FirstName.Default + " " + p.LastName.Default
	}

	for i := range roster.Forwards {
		wg.Add(1)
		go addFullName(&roster.Forwards[i])
	}
	for i := range roster.Defensemen {
		wg.Add(1)
		go addFullName(&roster.Defensemen[i])
	}
	for i := range roster.Goalies {
		wg.Add(1)
		go addFullName(&roster.Goalies[i])
	}

	wg.Wait()

	return &roster, nil
}
