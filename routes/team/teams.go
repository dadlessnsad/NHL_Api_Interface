package teams

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"nhl_interface/helpers"
	"nhl_interface/models"
	"nhl_interface/services"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

var db *sql.DB

func fetchTeamRoster(teamAbbr string) (*models.Team, error) {
	baseURL := os.Getenv("NHL_API_URL")
	apiURL := baseURL + "v1/roster/" + teamAbbr + "/current"
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmtRoster, err := helpers.FormatRosterData(teamAbbr, string(body))
	if err != nil {
		return nil, err
	}

	return fmtRoster, nil
}

func GetTeamRoster() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		TEAM_ABBR := mux.Vars(r)["TEAM_ABBR"]
		TEAM_ABBR = helpers.Capitalize(TEAM_ABBR)
		if TEAM_ABBR == "" {
			helpers.RespondWithError(w, http.StatusBadRequest, "Team abbreviation is required")
			return
		}
		if !helpers.VerifyTeamAbbr(TEAM_ABBR) {
			helpers.RespondWithError(w, http.StatusBadRequest, "Invalid team abbreviation")
			return
		}

		rosterData, err := fetchTeamRoster(TEAM_ABBR)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// insert team into database
		err = services.InsertTeam(*rosterData)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.RespondWithJSON(w, http.StatusOK, rosterData)
	}
}

func GetAllTeamsRosters() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		allAbbrs := helpers.GetAllTeamAbbrs()
		var wg sync.WaitGroup
		teamRosters := make(map[string]interface{})
		mutex := &sync.Mutex{}
		for _, abbr := range allAbbrs {
			wg.Add(1)
			go func(abbr string) {
				defer wg.Done()
				rosterData, err := fetchTeamRoster(abbr)
				mutex.Lock()
				defer mutex.Unlock()
				if err != nil {
					return
				}
				teamRosters[abbr] = rosterData
				if err != nil {
					return
				}
			}(abbr)
		}
		wg.Wait()
		helpers.RespondWithJSON(w, http.StatusOK, teamRosters)
	}
}

func GetAllTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		abbrs := helpers.GetAllTeamAbbrs()
		var teamsData []models.Team

		for _, abbr := range abbrs {
			team, err := services.QueryTeam(abbr)
			if err != nil {
				helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			teamsData = append(teamsData, team)
		}

		helpers.RespondWithJSON(w, http.StatusOK, teamsData)
	}
}
