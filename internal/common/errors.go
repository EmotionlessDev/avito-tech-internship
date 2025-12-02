package common

import (
	"errors"
	"net/http"

	"github.com/EmotionlessDev/avito-tech-internship/internal/helpers"
)

var ErrUserNotFound error = errors.New("user not found")
var ErrTeamNotFound error = errors.New("team not found")
var ErrTeamDuplicate error = errors.New("team already exists")
var ErrInvalidTeamName error = errors.New("invalid team name")
var ErrPRExists error = errors.New("pull request already exists")
var ErrPRNotFound error = errors.New("pull request not found")
var ErrAuthorNotFound error = errors.New("author not found")
var ErrNotAssignedReviewer error = errors.New("user is not an assigned reviewer for this pull request")
var ErrNoCandidates error = errors.New("no active replacement candidate in team")

func errorResponse(w http.ResponseWriter, status int, err error, message interface{}) {
	helpers.WriteJSON(w, status, helpers.Envelope{
		"error":   err.Error(),
		"message": message,
	}, nil)
}

func MethodNotAllowedResponse(w http.ResponseWriter) {
	errorResponse(w, http.StatusMethodNotAllowed, errors.New("method not allowed"), "the requested method is not allowed for the specified resource")
}

func InternalServerErrorResponse(w http.ResponseWriter, err error) {
	errorResponse(w, http.StatusInternalServerError, errors.New("internal server error"), err.Error())
}

func BadRequestResponse(w http.ResponseWriter, err error) {
	errorResponse(w, http.StatusBadRequest, errors.New("bad request"), err.Error())
}

func NotFoundResponse(w http.ResponseWriter, err error) {
	errorResponse(w, http.StatusNotFound, errors.New("not found"), err.Error())
}

func FailedValidationResponse(w http.ResponseWriter, errors map[string]string) {
	helpers.WriteJSON(w, http.StatusUnprocessableEntity, helpers.Envelope{
		"error":   "failed validation",
		"message": errors,
	}, nil)
}
