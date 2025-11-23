package http

import "net/http"

func MapTeamRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/team/add", h.AddTeam)
	mux.HandleFunc("/team/get", h.GetTeam)
}
