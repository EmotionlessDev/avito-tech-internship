package handlers

import (
	"log/slog"
	"net/http"

	"github.com/EmotionlessDev/avito-tech-internship/internal/helpers"
)

type ErrorResponder struct {
	Logger *slog.Logger
}

func NewErrorResponder(logger *slog.Logger) *ErrorResponder {
	return &ErrorResponder{
		Logger: logger,
	}
}

func (e *ErrorResponder) logError(r *http.Request, err error) {
	e.Logger.Error(err.Error(),
		slog.String("method", r.Method),
		slog.String("url", r.URL.String()))
}

func (e *ErrorResponder) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := helpers.Envelope{"error": message}
	err := helpers.WriteJSON(w, status, env, nil)
	if err != nil {
		e.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (e *ErrorResponder) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	e.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	e.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (e *ErrorResponder) NotFound(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	e.errorResponse(w, r, http.StatusNotFound, message)
}

func (e *ErrorResponder) BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	e.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (e *ErrorResponder) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	e.ServerError(w, r, err)
}

func (e *ErrorResponder) Conflict(w http.ResponseWriter, r *http.Request, message string) {
	e.errorResponse(w, r, http.StatusConflict, message)
}
