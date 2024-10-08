package main

import (
	"database/sql"
	"fmt"
	"icc-backend-test/database"
)

func GameEnded(playerWinner, playerLoser Player, draw bool, game *Game) {
	//* I assume that in the case of a draw, the ratings of both players remain the same
	if !draw {
		game.Winner = playerWinner
		game.Loser = playerLoser
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

	if rows.Next() {
		var game Game
		var whitePlayerID, blackPlayerID, winnerID, loserID int64

		if err := rows.Scan(&game.ID, &game.URL, &game.Rated, &game.GameType, &whitePlayerID, &blackPlayerID, &winnerID, &loserID, &game.Draw); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		game.WhitePlayer.ID = whitePlayerID
		game.BlackPlayer.ID = blackPlayerID
		game.Winner.ID = winnerID
		game.Loser.ID = loserID

		return &game, nil
	}
	return nil, fmt.Errorf("no game found with ID %d", id)
}

func CreateGame(playerWhite, playerBlack Player, gameType string, rated bool) *Game {
	game := &Game{
		WhitePlayer: playerWhite,
		BlackPlayer: playerBlack,
		GameType:    gameType,
		Rated:       rated,
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
