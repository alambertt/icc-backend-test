package model


type Game struct {
	ID            int64  `json:"id"`
	BlackPlayerID int64  `json:"black"`
	Draw          bool   `json:"draw"`
	GameType      string `json:"game_type"`
	LoserID       int64  `json:"loser_id"`
	URL           string `json:"url"`
	WhitePlayerID int64  `json:"white"`
	WinnerID      int64  `json:"winner_id"`
}