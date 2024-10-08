package main

var GameTypes []string = []string{"bullet", "blitz", "rapid", "classic"}

const ratingChange int64 = 10

type Player struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	BulletRating  int64  `json:"bullet_rating"`
	BlitzRating   int64  `json:"blitz_rating"`
	RapidRating   int64  `json:"rapid_rating"`
	ClassicRating int64  `json:"classic_rating"`
}

type Game struct {
	ID          int64  `json:"id"`
	BlackPlayer Player `json:"black"`
	Draw        bool   `json:"draw"`
	GameType    string `json:"game_type"`
	Loser       Player `json:"loser"`
	Rated       bool   `json:"rated"`
	URL         string `json:"url"`
	WhitePlayer Player `json:"white"`
	Winner      Player `json:"winner"`
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
