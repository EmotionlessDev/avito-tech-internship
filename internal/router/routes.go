package router

import (
	"github.com/EmotionlessDev/avito-tech-internship/internal/application"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(app *application.Application) *httprouter.Router {
	router := httprouter.New()

	// healthcheck
	router.HandlerFunc("GET", "/api/health", app.Handlers.Health.Check)

	return router
}
