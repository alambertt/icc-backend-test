package main

import (
	"database/sql"
	"fmt"
	"icc-backend-test/database"
	"math/rand"
	"time"
)

func PairingAlgorithm(player *Player, gameType string) {
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

	parsedRooms := parseRooms(rooms)
	if len(parsedRooms) == 0 {
		room := &Room{
			PlayerID:   player.ID,
			PlayerRate: playerRating,
			GameType:   gameType,
		}
		_, err := InsertRoomQuery(db, room)
		if err != nil {
			panic(err)
		}
	} else {
		randomIndex := getRandomNumber(0, len(parsedRooms)-1)
		room := parsedRooms[randomIndex]
		player2, err := FetchPlayerByID(db, room.PlayerID)
		if err != nil {
			panic(err)
		}
		game := CreateGame(*player, *player2, gameType, true)
		SendURLToPlayers(player, player2, game.URL)
	}
}

func getRandomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// SendURLToPlayers sends the URL of the game to the players using WebSockets. When the players receive the URL, they can start playing the game. The frontend will use the URL to redirect the players to the game.
func SendURLToPlayers(player1, player2 *Player, url string) {
	message := map[string]string{
		"type": "game_url",
		"url":  url,
	}

	// Send URL to white player
	if err := player1.Conn.WriteJSON(message); err != nil {
		fmt.Println("Error sending URL to player 1:", err)
	}

	// Send URL to black player
	if err := player2.Conn.WriteJSON(message); err != nil {
		fmt.Println("Error sending URL to player 2:", err)
	}
}

func parseRooms(rooms *sql.Rows) []Room {
	var parsedRooms []Room
	for rooms.Next() {
		var room Room
		if err := rooms.Scan(&room.ID, &room.PlayerID, &room.PlayerRate, &room.GameType); err != nil {
			panic(err)
		}
		parsedRooms = append(parsedRooms, room)
	}
	return parsedRooms
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
		URL:           createURL(playerWhite, playerBlack),
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

func createURL(playerWhite, playerBlack Player) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("https://play.chessclub.com/game/%s-vs-%s-%d", playerWhite.Name, playerBlack.Name, timestamp)
}
