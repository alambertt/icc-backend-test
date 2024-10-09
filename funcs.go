package main

import (
	"database/sql"
	"fmt"
	"icc-backend-test/constants"
	"icc-backend-test/database"
	"icc-backend-test/model"
	"icc-backend-test/utils"
	"icc-backend-test/websocket"
)

func PairingRequest(player *model.Player, gameType string) {
	var playerRating int64

	switch gameType {
	case "bullet":
		playerRating = player.BulletRating
	case "blitz":
		playerRating = player.BlitzRating
	case "rapid":
		playerRating = player.RapidRating
	case "classic":
		playerRating = player.ClassicRating
	}

	db, err := database.ConnectToMySQLDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rooms, err := GetRoomsByRateQuery(db, int(playerRating), gameType)
	if err != nil {
		panic(err)
	}
	defer rooms.Close()

	parsedRooms := utils.ParseRooms(rooms)
	if len(parsedRooms) == 0 {
		room := &model.Room{
			PlayerID:   player.ID,
			PlayerRate: playerRating,
			GameType:   gameType,
		}
		_, err := InsertRoomQuery(db, room)
		if err != nil {
			panic(err)
		}
	} else {
		randomIndex := utils.GetRandomNumber(0, len(parsedRooms)-1)
		room := parsedRooms[randomIndex]
		player2, err := FetchPlayerByID(db, room.PlayerID)
		if err != nil {
			panic(err)
		}
		game := CreateGame(*player, *player2, gameType, true)
		websocket.SendURLToPlayers(player, player2, game.URL)
		_, err = DeleteRoomsByPlayerIDQuery(db, int(player.ID))
		if err != nil {
			panic(err)
		}
		_, err = DeleteRoomsByPlayerIDQuery(db, int(player2.ID))
		if err != nil {
			panic(err)
		}
	}
}

func GameEnded(playerWinner, playerLoser model.Player, draw bool, game *model.Game) {
	//* I assume that in the case of a draw, the ratings of both players remain the same
	if !draw {
		game.WinnerID = playerWinner.ID
		game.LoserID = playerLoser.ID
	} else {
		game.Draw = true
	}
	db, err := database.ConnectToMySQLDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	UpdateGameQuery(db, game)

	playerWinner.ChangeRating(game.GameType, true)
	playerLoser.ChangeRating(game.GameType, false)
	UpdatePlayerQuery(db, &playerWinner)
	UpdatePlayerQuery(db, &playerLoser)
}

func CreateGame(playerWhite, playerBlack model.Player, gameType string, rated bool) *model.Game {
	if constants.GAME_TYPES[gameType] == "" {
		panic(fmt.Sprintf("invalid game type: %s", gameType))
	}
	game := &model.Game{
		WhitePlayerID: playerWhite.ID,
		BlackPlayerID: playerBlack.ID,
		GameType:      constants.GAME_TYPES[gameType],
		URL:           utils.CreateURL(playerWhite, playerBlack),
	}
	db, err := database.ConnectToMySQLDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := InsertGameQuery(db, game)
	if err != nil {
		panic(err)
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		panic(fmt.Sprintf("failed to get last insert ID: %v", err))
	}

	fetchedGame, err := FetchGameByID(db, lastInsertID)
	if err != nil {
		panic(err)
	}

	return fetchedGame
}

func CancelGameRequest(room *model.Room) {
	db, err := database.ConnectToMySQLDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = DeleteRoomQuery(db, room.GameType, int(room.PlayerID))
	if err != nil {
		panic(err)
	}
}

func FetchGameByID(db *sql.DB, id int64) (*model.Game, error) {
	rows, err := GetGameQuery(db, int(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %v", err)
	}
	defer rows.Close()

	var game model.Game
	if rows.Next() {
		if err := rows.Scan(&game.ID, &game.GameType, &game.WhitePlayerID, &game.BlackPlayerID, &game.WinnerID, &game.LoserID, &game.Draw); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		return &game, nil
	}
	return nil, fmt.Errorf("no game found with ID %d", id)
}

func FetchPlayerByID(db *sql.DB, id int64) (*model.Player, error) {
	rows, err := GetPlayerQuery(db, int(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get player: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		var player model.Player
		if err := rows.Scan(&player.ID, &player.Name, &player.BulletRating, &player.BlitzRating, &player.RapidRating, &player.ClassicRating); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		return &player, nil
	}
	return nil, fmt.Errorf("no player found with ID %d", id)
}
