package main

import (
	"database/sql"
	"fmt"
	"icc-backend-test/database"
)

func PairingAlgorithm(player *Player, gameType string) {
	// Implement a pairing algorithm
}

func GameEnded(playerWinner, playerLoser Player, draw bool, game *Game) {
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

func FetchGameByID(db *sql.DB, id int64) (*Game, error) {
	rows, err := GetGameQuery(db, int(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %v", err)
	}
	defer rows.Close()

	var game Game
	if rows.Next() {
		if err := rows.Scan(&game.ID, &game.GameType, &game.WhitePlayerID, &game.BlackPlayerID, &game.WinnerID, &game.LoserID, &game.Draw); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		return &game, nil
	}
	return nil, fmt.Errorf("no game found with ID %d", id)
}

func FetchPlayerByID(db *sql.DB, id int64) (*Player, error) {
	rows, err := GetPlayerQuery(db, int(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get player: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		var player Player
		if err := rows.Scan(&player.ID, &player.Name, &player.BulletRating, &player.BlitzRating, &player.RapidRating, &player.ClassicRating); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		return &player, nil
	}
	return nil, fmt.Errorf("no player found with ID %d", id)
}

func CreateGame(playerWhite, playerBlack Player, gameType string, rated bool) *Game {
	if GameTypes[gameType] == "" {
		panic(fmt.Sprintf("invalid game type: %s", gameType))
	}
	game := &Game{
		WhitePlayerID: playerWhite.ID,
		BlackPlayerID: playerBlack.ID,
		GameType:      GameTypes[gameType],
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
