package main

import (
	"database/sql"
	"fmt"
	"icc-backend-test/database"
)

func GetGameQuery(db *sql.DB, id int) (*sql.Rows, error) {
	query := "SELECT * FROM games WHERE id = ?"
	rows, err := database.ExecuteMySQLQuery(db, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %v", err)
	}
	return rows, nil
}

func GetPlayerQuery(db *sql.DB, id int) (*sql.Rows, error) {
	query := "SELECT * FROM WHERE id = ?"
	rows, err := database.ExecuteMySQLQuery(db, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player: %v", err)
	}
	return rows, nil
}

func InsertGameQuery(db *sql.DB, game *Game) (sql.Result, error) {
	query := "INSERT INTO games (id, url, rated, game_type, white_player_id, black_player_id, winner_id, loser_id, draw) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := database.ExecuteMySQLNonQuery(db, query, game.ID, game.URL, game.Rated, game.GameType, game.WhitePlayer.ID, game.BlackPlayer.ID, game.Winner.ID, game.Loser.ID, game.Draw)
	if err != nil {
		return nil, fmt.Errorf("failed to insert game: %v", err)
	}
	return result, nil
}

func InsertPlayerQuery(db *sql.DB, player *Player) (sql.Result, error) {
	query := "INSERT INTO (id, name, bullet_rating, blitz_rating, rapid_rating, classic_rating) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := database.ExecuteMySQLNonQuery(db, query, player.ID, player.Name, player.BulletRating, player.BlitzRating, player.RapidRating, player.ClassicRating)
	if err != nil {
		return nil, fmt.Errorf("failed to insert player: %v", err)
	}
	return result, nil
}

func UpdateGameQuery(db *sql.DB, game *Game) (sql.Result, error) {
	query := "UPDATE games SET url = ?, rated = ?, game_type = ?, white_player_id = ?, black_player_id = ?, winner_id = ?, loser_id = ?, draw = ? WHERE id = ?"
	result, err := database.ExecuteMySQLNonQuery(db, query, game.URL, game.Rated, game.GameType, game.WhitePlayer.ID, game.BlackPlayer.ID, game.Winner.ID, game.Loser.ID, game.Draw, game.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update game: %v", err)
	}
	return result, nil
}

func UpdatePlayerQuery(db *sql.DB, player *Player) (sql.Result, error) {
	query := "UPDATE SET name = ?, bullet_rating = ?, blitz_rating = ?, rapid_rating = ?, classic_rating = ? WHERE id = ?"
	result, err := database.ExecuteMySQLNonQuery(db, query, player.Name, player.BulletRating, player.BlitzRating, player.RapidRating, player.ClassicRating, player.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update player: %v", err)
	}
	return result, nil
}