package http

import (
	"context"
	"net/http"

	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/users"
	"github.com/EmotionlessDev/avito-tech-internship/internal/helpers"
)

type UpdateService interface {
	SetUserActiveByID(ctx context.Context, id string, isActive bool) (*users.User, error)
}

type Handler struct {
	updateService UpdateService
}

func NewHandler(updateService UpdateService) *Handler {
	return &Handler{
		updateService: updateService,
	}
}

func (h *Handler) SetUserActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.WriteJSON(w, http.StatusMethodNotAllowed, helpers.Envelope{"error": "method not allowed"}, nil)
		return
	}

	var req SetUserActiveRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"error": "invalid request"}, nil)
		return
	}

	ctx := context.Background()
	user, err := h.updateService.SetUserActiveByID(ctx, req.ID, req.IsActive)
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{"error": "failed to update user"}, nil)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"user": user}, nil)
}
