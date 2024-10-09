package main

import (
	"testing"

	"icc-backend-test/constants"
	"icc-backend-test/model"
	"icc-backend-test/utils"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPairingGameRequest_NoRooms(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	player := &model.Player{
		ID:            1,
		Name:          "Player1",
		BulletRating:  1500,
		BlitzRating:   1500,
		RapidRating:   1500,
		ClassicRating: 1500,
		Status:        "available",
	}
	gameType := "blitz"

	mock.ExpectQuery("SELECT \\* FROM rooms WHERE game_type = \\? AND player_rate BETWEEN \\? AND \\?").
		WithArgs(gameType, 1350, 1650).
		WillReturnRows(sqlmock.NewRows(nil))

	mock.ExpectExec("INSERT INTO rooms \\(player_id, player_rate, game_type\\) VALUES \\(\\?, \\?, \\?\\)").
		WithArgs(player.ID, player.BlitzRating, gameType).
		WillReturnResult(sqlmock.NewResult(1, 1))

	url, err := PairingGameRequest(db, player, gameType)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if url != "" {
		t.Errorf("expected empty URL, got %s", url)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPairingGameRequest_RoomsAvailable(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	player1 := &model.Player{
		ID:            1,
		Name:          "Player1",
		BulletRating:  1500,
		BlitzRating:   1500,
		RapidRating:   1500,
		ClassicRating: 1500,
		Status:        "available",
	}
	player2 := &model.Player{
		ID:            2,
		Name:          "Player2",
		BulletRating:  1500,
		BlitzRating:   1500,
		RapidRating:   1500,
		ClassicRating: 1500,
		Status:        "available",
	}
	gameType := "blitz"
	expectedGame := &model.Game{
		ID:            1,
		WhitePlayerID: player1.ID,
		BlackPlayerID: player2.ID,
		GameType:      constants.GAME_TYPES[gameType],
		URL:           utils.CreateURL(*player1, *player2),
	}

	mock.ExpectQuery("SELECT \\* FROM rooms WHERE game_type = \\? AND player_rate BETWEEN \\? AND \\?").
		WithArgs(gameType, 1350, 1650).
		WillReturnRows(sqlmock.NewRows([]string{"room_id", "player_id", "player_rate", "game_type"}).
			AddRow(1, player2.ID, player2.BlitzRating, gameType))

	mock.ExpectQuery("SELECT \\* FROM players WHERE id = \\?").
		WithArgs(player2.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "bullet_rating", "blitz_rating", "rapid_rating", "classic_rating", "status"}).
			AddRow(player2.ID, player2.Name, player2.BulletRating, player2.BlitzRating, player2.RapidRating, player2.ClassicRating, player2.Status))

	mock.ExpectExec("UPDATE players SET name = \\?, bullet_rating = \\?, blitz_rating = \\?, rapid_rating = \\?, classic_rating = \\?, status = \\? WHERE id = \\?").
		WithArgs(player1.Name, player1.BulletRating, player1.BlitzRating, player1.RapidRating, player1.ClassicRating, "playing", player1.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("UPDATE players SET name = \\?, bullet_rating = \\?, blitz_rating = \\?, rapid_rating = \\?, classic_rating = \\?, status = \\? WHERE id = \\?").
		WithArgs(player2.Name, player2.BulletRating, player2.BlitzRating, player2.RapidRating, player2.ClassicRating, "playing", player2.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO games \\(url, game_type, white_player_id, black_player_id, winner_id, loser_id, draw\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?\\)").
		WithArgs(expectedGame.URL, expectedGame.GameType, expectedGame.WhitePlayerID, expectedGame.BlackPlayerID, expectedGame.WinnerID, expectedGame.LoserID, expectedGame.Draw).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT \\* FROM games WHERE id = \\?").
		WithArgs(expectedGame.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "url", "game_type", "white_player_id", "black_player_id", "winner_id", "loser_id", "draw"}).
			AddRow(expectedGame.ID, expectedGame.URL, expectedGame.GameType, expectedGame.WhitePlayerID, expectedGame.BlackPlayerID, expectedGame.WinnerID, expectedGame.LoserID, expectedGame.Draw))

	mock.ExpectExec("DELETE FROM rooms WHERE player_id = \\?").
		WithArgs(player1.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM rooms WHERE player_id = \\?").
		WithArgs(player2.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	url, err := PairingGameRequest(db, player1, gameType)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if url != expectedGame.URL {
		t.Errorf("expected URL %s, got %s", expectedGame.URL, url)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
