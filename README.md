Leaderboard API

A RESTful leaderboard service built with Go, Redis, and Prometheus.
Supports score submission, rank retrieval, pagination, and user removal with performance monitoring.

Features

Submit user scores

Get paginated leaderboard

Fetch user rank and score

Remove user from leaderboard

Prometheus metrics integration

Redis Sorted Sets for efficient ranking

Tech Stack

Go

Redis

Chi Router

Prometheus

API Endpoints
POST   /v1/leaderboards/{boardId}/scores
GET    /v1/leaderboards/{boardId}
GET    /v1/leaderboards/{boardId}/users/{userId}
DELETE /v1/leaderboards/{boardId}/users/{userId}
GET    /metrics
Run Locally
go run main.go

Ensure Redis is running on localhost:6379 or configure via environment variables.

Run Tests
go test ./... -v
