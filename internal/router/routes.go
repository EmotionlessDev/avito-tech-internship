package router

import (
	"github.com/EmotionlessDev/avito-tech-internship/internal/application"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(app *application.Application) *httprouter.Router {
	router := httprouter.New()

	// healthcheck
	router.HandlerFunc("GET", "/health", app.Handlers.Health.Check)

	// team
	router.HandlerFunc("POST", "/team/add", app.Handlers.Team.CreateTeam)
	router.HandlerFunc("GET", "/team/get", app.Handlers.Team.GetTeam)

	// user
	router.HandlerFunc("POST", "/users/setIsActive", app.Handlers.User.SetIsActive)

	return router
}
