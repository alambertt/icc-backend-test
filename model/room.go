package model

type Room struct {
	ID         int64  `json:"id"`
	PlayerID   int64  `json:"player_id"`
	PlayerRate int64  `json:"player_rate"`
	GameType   string `json:"game_type"`
}
