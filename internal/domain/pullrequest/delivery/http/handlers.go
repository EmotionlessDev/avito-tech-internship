package http

import (
	"context"
	"net/http"

	"github.com/EmotionlessDev/avito-tech-internship/internal/common"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/pullrequest"
	"github.com/EmotionlessDev/avito-tech-internship/internal/helpers"
)

type CreateService interface {
	CreatePR(ctx context.Context, pr pullrequest.PullRequest) (*pullrequest.PullRequestWithReviewers, error)
}

type Handler struct {
	service CreateService
}

func NewHandler(service CreateService) *Handler {
	return &Handler{service: service}
}

type CreatePRRequest struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

func (h *Handler) CreatePR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.WriteJSON(w, http.StatusMethodNotAllowed, helpers.Envelope{"error": "method not allowed"}, nil)
		return
	}

	var req CreatePRRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"error": "invalid_request", "message": err.Error()}, nil)
		return
	}

	pr := pullrequest.PullRequest{
		ID:       req.ID,
		Name:     req.Name,
		AuthorID: req.AuthorID,
	}

	createdPR, err := h.service.CreatePR(r.Context(), pr)
	if err != nil {
		switch err {
		case common.ErrPRExists:
			helpers.WriteJSON(w, http.StatusConflict, helpers.Envelope{"error": "PR_EXISTS", "message": "PR id already exists"}, nil)
		default:
			helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{"error": "internal_error", "message": err.Error()}, nil)
		}
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, helpers.Envelope{"pr": createdPR}, nil)
}
