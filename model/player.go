package model

import (
	"icc-backend-test/constants"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	BulletRating  int64  `json:"bullet_rating"`
	BlitzRating   int64  `json:"blitz_rating"`
	RapidRating   int64  `json:"rapid_rating"`
	ClassicRating int64  `json:"classic_rating"`
	Conn          *websocket.Conn
}

func (p Player) ChangeRating(gameType string, up bool) {
	switch gameType {
	case "bullet":
		if up {
			p.BulletRating = p.BulletRating + constants.RATING_CHANGE
		}
		p.BulletRating = p.BulletRating - constants.RATING_CHANGE
	case "blitz":
		if up {
			p.BlitzRating = p.BlitzRating + constants.RATING_CHANGE
		}
		p.BlitzRating = p.BlitzRating - constants.RATING_CHANGE
	case "rapid":
		if up {
			p.RapidRating = p.RapidRating + constants.RATING_CHANGE
		}
		p.RapidRating = p.RapidRating - constants.RATING_CHANGE
	case "classic":
		if up {
			p.ClassicRating = p.ClassicRating + constants.RATING_CHANGE
		}
		p.ClassicRating = p.ClassicRating - constants.RATING_CHANGE
	}
}
