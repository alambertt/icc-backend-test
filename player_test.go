package main

import (
	"testing"

	"icc-backend-test/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateNewPlayer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO players \\(name, bullet_rating, blitz_rating, rapid_rating, classic_rating, status\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?\\)").
		WithArgs("test_player", 1500, 1500, 1500, 1500, "available").
		WillReturnResult(sqlmock.NewResult(1, 1))

	player, err := CreateNewPlayer(db, "test_player")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if player.BulletRating != 1500 || player.BlitzRating != 1500 || player.RapidRating != 1500 || player.ClassicRating != 1500 {
		t.Errorf("expected initial ratings to be 1500, got BulletRating: %d, BlitzRating: %d, RapidRating: %d, ClassicRating: %d",
			player.BulletRating, player.BlitzRating, player.RapidRating, player.ClassicRating)
	}

	if player.Status != "available" {
		t.Errorf("expected status to be 'available', got %s", player.Status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFetchPlayerByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	playerID := int64(1)
	expectedPlayer := model.Player{
		ID:            playerID,
		Name:          "test_player",
		BulletRating:  1500,
		BlitzRating:   1500,
		RapidRating:   1500,
		ClassicRating: 1500,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "bullet_rating", "blitz_rating", "rapid_rating", "classic_rating", "status"}).
		AddRow(expectedPlayer.ID, expectedPlayer.Name, expectedPlayer.BulletRating, expectedPlayer.BlitzRating, expectedPlayer.RapidRating, expectedPlayer.ClassicRating, "available")

	mock.ExpectQuery("SELECT \\* FROM players WHERE id = \\?").
		WithArgs(playerID).
		WillReturnRows(rows)

	player, err := FetchPlayerByID(db, playerID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if player.ID != expectedPlayer.ID || player.Name != expectedPlayer.Name || player.BulletRating != expectedPlayer.BulletRating ||
		player.BlitzRating != expectedPlayer.BlitzRating || player.RapidRating != expectedPlayer.RapidRating || player.ClassicRating != expectedPlayer.ClassicRating {
		t.Errorf("expected player %+v, got %+v", expectedPlayer, player)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

