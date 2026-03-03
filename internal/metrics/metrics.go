package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ScoreSubmitted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "leaderboard_score_submitted_total",
		Help: "Total number of score submissions",
	})

	LeaderboardRead = promauto.NewCounter(prometheus.CounterOpts{
		Name: "leaderboard_read_total",
		Help: "Total number of leaderboard read requests",
	})

	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "leaderboard_request_duration_seconds",
		Help:    "Request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"endpoint", "method", "status"})
)
