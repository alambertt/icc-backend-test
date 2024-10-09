package utils

import (
	"database/sql"
	"fmt"
	"icc-backend-test/model"
	"time"

	"math/rand"
)

func CreateURL(playerWhite, playerBlack model.Player) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("https://play.chessclub.com/game/%s-vs-%s-%d", playerWhite.Name, playerBlack.Name, timestamp)
}

func GetRandomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func ParseRooms(rooms *sql.Rows) []model.Room {
	var parsedRooms []model.Room
	for rooms.Next() {
		var room model.Room
		if err := rooms.Scan(&room.ID, &room.PlayerID, &room.PlayerRate, &room.GameType); err != nil {
			panic(err)
		}
		parsedRooms = append(parsedRooms, room)
	}
	return parsedRooms
}