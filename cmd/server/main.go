package main

import (
	"log"
	"net/http"
	"os"

	"leaderBoard/internal/api"
	"leaderBoard/internal/config"
	"leaderBoard/internal/service"
	"leaderBoard/internal/store/redis"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.LoadConfig()
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisClient := config.InitRedis(cfg)
	redisStore := redis.NewRedisStore(redisClient)
	svc := service.NewService(redisStore)
	handler := api.NewHandler(svc)
	r := chi.NewRouter()

	r.Route("/v1/leaderboards/{boardId}", func(r chi.Router) {
		r.Post("/scores", handler.SubmitScore)
		r.Get("/", handler.GetLeaderboard)
		r.Get("/users/{userId}", handler.GetUserRank)
		r.Delete("/users/{userId}", handler.RemoveUser)
	})
	r.Handle("/metrics", promhttp.Handler())

	log.Println("Server running on: ", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
