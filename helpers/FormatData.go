package helpers

import (
	"encoding/json"
	"fmt"
	"nhl_interface/models"
	"sync"
)

func FormatRosterData(teamAbbr string, rosterData string) (*models.Team, error) {
	var team models.Team
	var roster models.RosterResponse
	err := json.Unmarshal([]byte(rosterData), &roster)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling roster data: %v", err)
	}

	// get team name from team abbreviation
	teamName := GetTeamNameFromAbbr(teamAbbr)
	team.Name = teamName
	team.Abbrev = teamAbbr

	var wg sync.WaitGroup
	addFullName := func(p *models.Player) {
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

	team.Roster = roster
	//create a new team struct
	_team := models.Team{
		Name:   teamName,
		Abbrev: teamAbbr,
		Roster: roster,
	}

	return &_team, nil

}
