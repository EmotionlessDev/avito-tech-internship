package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/EmotionlessDev/avito-tech-internship/internal/helpers"
	"github.com/EmotionlessDev/avito-tech-internship/internal/repository"
	"github.com/EmotionlessDev/avito-tech-internship/internal/services"
)

type UserHandler struct {
	userService    services.UserService
	ErrorResponder *ErrorResponder
	Logger         *slog.Logger
}

func NewUserHandler(
	userService services.UserService,
	ErrorResponder *ErrorResponder,
	logger *slog.Logger,
) *UserHandler {
	return &UserHandler{
		userService:    userService,
		ErrorResponder: ErrorResponder,
		Logger:         logger,
	}
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := helpers.ReadJSON(w, r, &input); err != nil {
		h.ErrorResponder.BadRequest(w, r, err)
		return
	}

	user, err := h.userService.SetIsActive(r.Context(), input.UserID, input.IsActive)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			h.ErrorResponder.NotFound(w, r)
			return
		}
		h.ErrorResponder.InternalServerError(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"user": user}, nil)
	if err != nil {
		h.ErrorResponder.ServerError(w, r, err)
		return
	}
}
