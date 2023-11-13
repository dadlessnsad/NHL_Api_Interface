package routes

import (
	"net/http"
	teams "nhl_interface/routes/team"

	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//_logger.Info(`Incoming request to "` + req.RequestURI + `" from ` + req.RemoteAddr)
		next.ServeHTTP(w, req)
	})
}

func HandleRoutes(r *mux.Router) {
	r.Use(loggingMiddleware)

	r.HandleFunc(`/api/team/all_teams`, teams.GetAllTeams()).Methods("GET")

	r.HandleFunc(`/api/team/all_teams/roster`, teams.GetAllTeamsRosters()).Methods("GET")

	r.HandleFunc(`/api/team/{TEAM_ABBR}/roster`, teams.GetTeamRoster()).Methods("GET")

}
