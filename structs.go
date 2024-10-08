package main

import "github.com/gorilla/websocket"

var GameTypes map[string]string = map[string]string{
	"bullet":  "bullet",
	"blitz":   "blitz",
	"rapid":   "rapid",
	"classic": "classic",
}

const ratingChange int64 = 10
const roomRateEligibility int = 150

type Player struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	BulletRating  int64  `json:"bullet_rating"`
	BlitzRating   int64  `json:"blitz_rating"`
	RapidRating   int64  `json:"rapid_rating"`
	ClassicRating int64  `json:"classic_rating"`
	Conn 		*websocket.Conn
}

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

type Room struct {
	ID         int64  `json:"id"`
	PlayerID   int64  `json:"player_id"`
	PlayerRate int64  `json:"player_rate"`
	GameType   string `json:"game_type"`
}

func (p Player) ChangeRating(gameType string, up bool) {
	switch gameType {
	case "bullet":
		if up {
			p.BulletRating = p.BulletRating + ratingChange
		}
		p.BulletRating = p.BulletRating - ratingChange
	case "blitz":
		if up {
			p.BlitzRating = p.BlitzRating + ratingChange
		}
		p.BlitzRating = p.BlitzRating - ratingChange
	case "rapid":
		if up {
			p.RapidRating = p.RapidRating + ratingChange
		}
		p.RapidRating = p.RapidRating - ratingChange
	case "classic":
		if up {
			p.ClassicRating = p.ClassicRating + ratingChange
		}
		p.ClassicRating = p.ClassicRating - ratingChange
	}
}
