package snapi

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHandlers(r *mux.Router)  {
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	r.HandleFunc(`/incidents/{incident:INC\d{7,10}}`, IncidentHandler) //get details for single incident
	r.HandleFunc("/incidents/{option:count}/{team}", IncidentTeamHandler) //get count of active incidents for team
	r.HandleFunc("/incidents/{option:list}/{team}", IncidentTeamHandler) //get list of active incidents for team
}
