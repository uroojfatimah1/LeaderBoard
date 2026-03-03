package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"leaderBoard/internal/api"
	"leaderBoard/internal/config"
	"leaderBoard/internal/service"
	"leaderBoard/internal/store/redis"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func setupTestHandler(t *testing.T) *api.Handler {
	cfg := config.LoadConfig()
	rdb := config.InitRedis(cfg)
	redisStore := redis.NewRedisStore(rdb)
	svc := service.NewService(redisStore)
	return api.NewHandler(svc)
}

// helper to create router for testing
func setupRouter(h *api.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/v1/leaderboards/{boardId}", func(r chi.Router) {
		r.Post("/scores", h.SubmitScore)
		r.Get("/", h.GetLeaderboard)
		r.Get("/users/{userId}", h.GetUserRank)
		r.Delete("/users/{userId}", h.RemoveUser)
	})
	return r
}

func TestLeaderboardCRUD(t *testing.T) {
	h := setupTestHandler(t)
	r := setupRouter(h)

	boardID := "test-board"
	userID := "user111"

	// --- Test SubmitScore ---
	scoreReq := map[string]interface{}{
		"userId": userID,
		"score":  1000,
	}
	body, _ := json.Marshal(scoreReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/leaderboards/"+boardID+"/scores", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, 1000.0, resp["score"])
	require.EqualValues(t, 1, resp["rank"])

	// --- Test GetLeaderboard ---
	req = httptest.NewRequest(http.MethodGet, "/v1/leaderboards/"+boardID+"?limit=10&offset=0", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var lbResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &lbResp)
	require.NoError(t, err)
	require.Equal(t, boardID, lbResp["boardId"])
	items := lbResp["items"].([]interface{})
	require.Len(t, items, 1)

	// --- Test GetUserRank ---
	req = httptest.NewRequest(http.MethodGet, "/v1/leaderboards/"+boardID+"/users/"+userID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var userResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &userResp)
	require.NoError(t, err)
	require.Equal(t, 1000.0, userResp["score"])
	require.EqualValues(t, 1, userResp["rank"])

	// --- Test RemoveUser ---
	req = httptest.NewRequest(http.MethodDelete, "/v1/leaderboards/"+boardID+"/users/"+userID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)

	// Verify user is removed
	req = httptest.NewRequest(http.MethodGet, "/v1/leaderboards/"+boardID+"/users/"+userID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}
