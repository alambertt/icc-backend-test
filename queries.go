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

func GetRoomsByRateQuery(db *sql.DB, rate int, gameType string) (*sql.Rows, error) {
	query := "SELECT * FROM rooms WHERE game_type = ? AND player_rate >= ? AND player_rate <= ?"
	rows, err := database.ExecuteMySQLQuery(db, query, gameType, rate-roomRateEligibility, rate+roomRateEligibility)
	if err != nil {
		return nil, fmt.Errorf("failed to get player by rate: %v", err)
	}
	return rows, nil
}

func GetGameTypeParameter(gameType string) (string, error) {
	if _, ok := GameTypes[gameType]; !ok {
		return "", fmt.Errorf("invalid game type: %s", gameType)
	}
	return GameTypes[gameType], nil
}

func InsertGameQuery(db *sql.DB, game *Game) (sql.Result, error) {
	query := "INSERT INTO games (id, url, game_type, white_player_id, black_player_id, winner_id, loser_id, draw) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := database.ExecuteMySQLNonQuery(db, query, game.ID, game.URL, game.GameType, game.WhitePlayerID, game.BlackPlayerID, game.WinnerID, game.LoserID, game.Draw)
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

func InsertRoomQuery(db *sql.DB, room *Room) (sql.Result, error) {
	query := "INSERT INTO rooms (id, player_id, player_rate, game_type) VALUES (?, ?, ?, ?)"
	result, err := database.ExecuteMySQLNonQuery(db, query, room.ID, room.PlayerID, room.PlayerRate, room.GameType)
	if err != nil {
		return nil, fmt.Errorf("failed to insert room: %v", err)
	}
	return result, nil
}

func UpdateGameQuery(db *sql.DB, game *Game) (sql.Result, error) {
	query := "UPDATE games SET game_type = ?, url = ?, white_player_id = ?, black_player_id = ?, winner_id = ?, loser_id = ?, draw = ? WHERE id = ?"
	result, err := database.ExecuteMySQLNonQuery(db, query, game.GameType,game.URL, game.WhitePlayerID, game.BlackPlayerID, game.WinnerID, game.LoserID, game.Draw, game.ID)
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
