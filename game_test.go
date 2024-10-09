package main

import (
	"fmt"
	"testing"

	"icc-backend-test/constants"
	"icc-backend-test/model"
	"icc-backend-test/utils"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFetchGameByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gameID := int64(1)
	expectedGame := model.Game{
		ID:            gameID,
		GameType:      "blitz",
		WhitePlayerID: 1,
		BlackPlayerID: 2,
		WinnerID:      1,
		LoserID:       2,
		Draw:          false,
	}

	rows := sqlmock.NewRows([]string{"id", "url", "game_type", "white_player_id", "black_player_id", "winner_id", "loser_id", "draw"}).
		AddRow(expectedGame.ID, expectedGame.URL, expectedGame.GameType, expectedGame.WhitePlayerID, expectedGame.BlackPlayerID, expectedGame.WinnerID, expectedGame.LoserID, expectedGame.Draw)

	mock.ExpectQuery("SELECT \\* FROM games WHERE id = \\?").
		WithArgs(gameID).
		WillReturnRows(rows)

	game, err := FetchGameByID(db, gameID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if game.ID != expectedGame.ID || game.GameType != expectedGame.GameType || game.WhitePlayerID != expectedGame.WhitePlayerID ||
		game.BlackPlayerID != expectedGame.BlackPlayerID || game.WinnerID != expectedGame.WinnerID || game.LoserID != expectedGame.LoserID || game.Draw != expectedGame.Draw {
		t.Errorf("expected game %+v, got %+v", expectedGame, game)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateGame(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	playerWhite := model.Player{
		ID:            1,
		Name:          "PlayerWhite",
		BulletRating:  1500,
		BlitzRating:   1500,
		RapidRating:   1500,
		ClassicRating: 1500,
		Status:        "available",
	}
	playerBlack := model.Player{
		ID:            2,
		Name:          "PlayerBlack",
		BulletRating:  1500,
		BlitzRating:   1500,
		RapidRating:   1500,
		ClassicRating: 1500,
		Status:        "available",
	}
	gameType := "blitz"
	expectedGame := &model.Game{
		ID:            1,
		WhitePlayerID: playerWhite.ID,
		BlackPlayerID: playerBlack.ID,
		GameType:      constants.GAME_TYPES[gameType],
		URL:           utils.CreateURL(playerWhite, playerBlack),
	}

	mock.ExpectExec("UPDATE players SET name = \\?, bullet_rating = \\?, blitz_rating = \\?, rapid_rating = \\?, classic_rating = \\?, status = \\? WHERE id = \\?").
		WithArgs(playerWhite.Name, playerWhite.BulletRating, playerWhite.BlitzRating, playerWhite.RapidRating, playerWhite.ClassicRating, "playing", playerWhite.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("UPDATE players SET name = \\?, bullet_rating = \\?, blitz_rating = \\?, rapid_rating = \\?, classic_rating = \\?, status = \\? WHERE id = \\?").
		WithArgs(playerBlack.Name, playerBlack.BulletRating, playerBlack.BlitzRating, playerBlack.RapidRating, playerBlack.ClassicRating, "playing", playerBlack.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO games \\(url, game_type, white_player_id, black_player_id, winner_id, loser_id, draw\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?\\)").
		WithArgs(expectedGame.URL, expectedGame.GameType, expectedGame.WhitePlayerID, expectedGame.BlackPlayerID, expectedGame.WinnerID, expectedGame.LoserID, expectedGame.Draw).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"id", "url", "game_type", "white_player_id", "black_player_id", "winner_id", "loser_id", "draw"}).
		AddRow(expectedGame.ID, expectedGame.URL, expectedGame.GameType, expectedGame.WhitePlayerID, expectedGame.BlackPlayerID, expectedGame.WinnerID, expectedGame.LoserID, expectedGame.Draw)

	mock.ExpectQuery("SELECT \\* FROM games WHERE id = \\?").
		WithArgs(expectedGame.ID).
		WillReturnRows(rows)

	game, err := CreateGame(db, playerWhite, playerBlack, gameType, true)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if game.ID != expectedGame.ID || game.GameType != expectedGame.GameType || game.WhitePlayerID != expectedGame.WhitePlayerID ||
		game.BlackPlayerID != expectedGame.BlackPlayerID || game.URL != expectedGame.URL {
		t.Errorf("expected game %+v, got %+v", expectedGame, game)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGameEndedWithWinnerAndLoser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	playerWinner := model.Player{
		ID:            1,
		Name:          "PlayerWinner",
		BulletRating:  1500,
		BlitzRating:   1500,
		RapidRating:   1500,
		ClassicRating: 1500,
		Status:        "playing",
	}
	playerLoser := model.Player{
		ID:            2,
		Name:          "PlayerLoser",
		BulletRating:  1500,
		BlitzRating:   1500,
		RapidRating:   1500,
		ClassicRating: 1500,
		Status:        "playing",
	}
	game := &model.Game{
		ID:            1,
		WhitePlayerID: playerWinner.ID,
		BlackPlayerID: playerLoser.ID,
		GameType:      "blitz",
	}

	mock.ExpectExec("UPDATE games SET game_type = \\?, url = \\?, white_player_id = \\?, black_player_id = \\?, winner_id = \\?, loser_id = \\?, draw = \\? WHERE id = \\?").
		WithArgs(game.GameType, game.URL, game.WhitePlayerID, game.BlackPlayerID, playerWinner.ID, playerLoser.ID, false, game.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("UPDATE players SET name = \\?, bullet_rating = \\?, blitz_rating = \\?, rapid_rating = \\?, classic_rating = \\?, status = \\? WHERE id = \\?").
		WithArgs(playerWinner.Name, playerWinner.BulletRating, 1510, playerWinner.RapidRating, playerWinner.ClassicRating, "available", playerWinner.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("UPDATE players SET name = \\?, bullet_rating = \\?, blitz_rating = \\?, rapid_rating = \\?, classic_rating = \\?, status = \\? WHERE id = \\?").
		WithArgs(playerLoser.Name, playerLoser.BulletRating, 1490, playerLoser.RapidRating, playerLoser.ClassicRating, "available", playerLoser.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = GameEnded(db, &playerWinner, &playerLoser, false, game)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGameEndedWithDraw(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	player1 := model.Player{
		ID:            1,
		Name:          "Player1",
		BulletRating:  1500,
		BlitzRating:   1500,
		RapidRating:   1500,
		ClassicRating: 1500,
		Status:        "playing",
	}
	player2 := model.Player{
		ID:            2,
		Name:          "Player2",
		BulletRating:  1500,
		BlitzRating:   1500,
		RapidRating:   1500,
		ClassicRating: 1500,
		Status:        "playing",
	}
	game := &model.Game{
		ID:            1,
		WhitePlayerID: player1.ID,
		BlackPlayerID: player2.ID,
		URL:           "https://play.chessclub.com/game/1",
		GameType:      "blitz",
	}

	mock.ExpectExec("UPDATE games SET game_type = \\?, url = \\?, white_player_id = \\?, black_player_id = \\?, winner_id = \\?, loser_id = \\?, draw = \\? WHERE id = \\?").
		WithArgs(game.GameType, game.URL, game.WhitePlayerID, game.BlackPlayerID, 0, 0, true, game.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("UPDATE players SET name = \\?, bullet_rating = \\?, blitz_rating = \\?, rapid_rating = \\?, classic_rating = \\?, status = \\? WHERE id = \\?").
		WithArgs(player1.Name, player1.BulletRating, player1.BlitzRating, player1.RapidRating, player1.ClassicRating, "available", player1.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("UPDATE players SET name = \\?, bullet_rating = \\?, blitz_rating = \\?, rapid_rating = \\?, classic_rating = \\?, status = \\? WHERE id = \\?").
		WithArgs(player2.Name, player2.BulletRating, player2.BlitzRating, player2.RapidRating, player2.ClassicRating, "available", player2.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = GameEnded(db, &player1, &player2, true, game)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCancelGameRequest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("successful room deletion", func(t *testing.T) {
		room := &model.Room{
			PlayerID: 1,
			GameType: "blitz",
		}

		mock.ExpectExec("DELETE FROM rooms WHERE player_id = \\? AND game_type = \\?").
			WithArgs(room.PlayerID, room.GameType).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := CancelGameRequest(db, room)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("failed room deletion", func(t *testing.T) {
		room := &model.Room{
			PlayerID: 1,
			GameType: "blitz",
		}

		mock.ExpectExec("DELETE FROM rooms WHERE player_id = \\? AND game_type = \\?").
			WithArgs( room.PlayerID, room.GameType).
			WillReturnError(fmt.Errorf("some error"))

		err := CancelGameRequest(db, room)
		if err == nil {
			t.Errorf("expected an error but got none")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
