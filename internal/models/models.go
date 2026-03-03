package models

type SubmitScoreRequest struct {
	UserID string  `json:"userId"`
	Score  float64 `json:"score"`
}

type UserRankResponse struct {
	Rank  int64   `json:"rank"`
	Score float64 `json:"score"`
}

type LeaderboardItem struct {
	UserID string  `json:"userId"`
	Rank   int64   `json:"rank"`
	Score  float64 `json:"score"`
}

type LeaderboardPage struct {
	BoardID string            `json:"boardId"`
	Total   int64             `json:"total"`
	Items   []LeaderboardItem `json:"items"`
}
