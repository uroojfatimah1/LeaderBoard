package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"leaderBoard/internal/metrics"
	"leaderBoard/internal/models"
	"leaderBoard/internal/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) SubmitScore(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	boardID := chi.URLParam(r, "boardId")

	var req models.SubmitScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	if req.Score < 0 {
		http.Error(w, "score must be non-negative", http.StatusBadRequest)
		return
	}

	rank, score, err := h.service.SubmitScore(r.Context(), boardID, req.UserID, req.Score)
	if err != nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}

	response := models.UserRankResponse{
		Rank:  rank,
		Score: score,
	}
	duration := time.Since(startTime).Seconds()
	metrics.RequestDuration.WithLabelValues("/scores", "POST", "200").Observe(duration)
	metrics.ScoreSubmitted.Inc()

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	boardID := chi.URLParam(r, "boardId")

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := int64(50)
	offset := int64(0)

	if limitStr != "" {
		val, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil || val < 0 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = val
	}

	if offsetStr != "" {
		val, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil || val < 0 {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}
		offset = val
	}

	entries, total, err := h.service.GetLeaderboard(r.Context(), boardID, offset, limit)
	if err != nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}

	items := make([]models.LeaderboardItem, len(entries))
	for i, e := range entries {
		items[i] = models.LeaderboardItem{
			UserID: e.UserID,
			Rank:   offset + int64(i) + 1,
			Score:  e.Score,
		}
	}

	response := models.LeaderboardPage{
		BoardID: boardID,
		Total:   total,
		Items:   items,
	}
	duration := time.Since(startTime).Seconds()
	metrics.RequestDuration.WithLabelValues("/leaderboard", "GET", "200").Observe(duration)
	metrics.LeaderboardRead.Inc()

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetUserRank(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	boardID := chi.URLParam(r, "boardId")
	userID := chi.URLParam(r, "userId")

	rank, score, err := h.service.GetUserRank(r.Context(), boardID, userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	response := models.UserRankResponse{
		Rank:  rank,
		Score: score,
	}
	duration := time.Since(startTime).Seconds()
	metrics.RequestDuration.WithLabelValues("/scoreById", "GET", "200").Observe(duration)
	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) RemoveUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	boardID := chi.URLParam(r, "boardId")
	userID := chi.URLParam(r, "userId")

	err := h.service.RemoveUser(r.Context(), boardID, userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}
	duration := time.Since(startTime).Seconds()
	metrics.RequestDuration.WithLabelValues("/scoreById", "DELETE", "204").Observe(duration)
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
