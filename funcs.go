package main

import (
	"database/sql"
	"fmt"
	"icc-backend-test/constants"
	"icc-backend-test/model"
	"icc-backend-test/utils"
	"icc-backend-test/websocket"
)

func PairingGameRequest(db *sql.DB, player *model.Player, gameType string) (string, error) {
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

	rooms, err := GetRoomsByRateQuery(db, int(playerRating), gameType)
	if err != nil {
		return "", err
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
			return "", err
		}
		return "", nil //! Return empty string because the player is waiting for an opponent. The frontend will display a message to the player telling him that he is waiting for an opponent.
	} else {
		randomIndex := utils.GetRandomNumber(0, len(parsedRooms)-1)
		room := parsedRooms[randomIndex]
		player2, err := FetchPlayerByID(db, room.PlayerID)
		if err != nil {
			return "", err
		}

		for player2.IsPlaying() {
			parsedRooms = utils.DeleteRoomFromArray(parsedRooms, randomIndex)
			randomIndex = utils.GetRandomNumber(0, len(parsedRooms)-1)
			room = parsedRooms[randomIndex]
			player2, err = FetchPlayerByID(db, room.PlayerID)
			if err != nil {
				return "", err
			}
		}
		game, err := CreateGame(db, *player, *player2, gameType, true)
		if err != nil {
			return "", err
		}
		// Send URL to player 2 using WebSockets because player 1 is the one who requested the game
		websocket.SendURLToPlayer(player2, game.URL)
		_, err = DeleteRoomsByPlayerIDQuery(db, int(player.ID))
		if err != nil {
			return "", err
		}
		_, err = DeleteRoomsByPlayerIDQuery(db, int(player2.ID))
		if err != nil {
			return "", err
		}
		return game.URL, nil
	}
}

func CancelGameRequest(db *sql.DB, room *model.Room) error {
	_, err := DeleteRoomQuery(db, room.GameType, int(room.PlayerID))
	return err
}

func GameEnded(db *sql.DB, playerWinner, playerLoser model.Player, draw bool, game *model.Game) error {
	//* I assume that in the case of a draw, the ratings of both players remain the same
	if !draw {
		game.WinnerID = playerWinner.ID
		game.LoserID = playerLoser.ID
		playerWinner.ChangeRating(game.GameType, true)
		playerLoser.ChangeRating(game.GameType, false)
	} else {
		game.Draw = true
	}

	_, err := UpdateGameQuery(db, game)
	if err != nil {
		return err
	}
	playerWinner.SetAvailable()
	_, err = UpdatePlayerQuery(db, &playerWinner)
	if err != nil {
		return err
	}

	playerLoser.SetAvailable()
	_, err = UpdatePlayerQuery(db, &playerLoser)
	if err != nil {
		return err
	}
	return nil
}

func CreateGame(db *sql.DB, playerWhite, playerBlack model.Player, gameType string, rated bool) (*model.Game, error) {
	if constants.GAME_TYPES[gameType] == "" {
		return nil, fmt.Errorf("invalid game type: %s", gameType)
	}
	game := &model.Game{
		WhitePlayerID: playerWhite.ID,
		BlackPlayerID: playerBlack.ID,
		GameType:      constants.GAME_TYPES[gameType],
		URL:           utils.CreateURL(playerWhite, playerBlack),
	}

	playerWhite.SetPlaying()
	playerBlack.SetPlaying()
	_, err := UpdatePlayerQuery(db, &playerWhite)
	if err != nil {
		return nil, err
	}
	_, err = UpdatePlayerQuery(db, &playerBlack)
	if err != nil {
		return nil, err
	}

	res, err := InsertGameQuery(db, game)
	if err != nil {
		return nil, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		panic(fmt.Sprintf("failed to get last insert ID: %v", err))
	}

	fetchedGame, err := FetchGameByID(db, lastInsertID)
	if err != nil {
		return nil, err
	}

	return fetchedGame, nil
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

func CreateNewPlayer(db *sql.DB, name string) (*model.Player, error) {
	player := &model.Player{
		Name:          name,
		BulletRating:  1500,
		BlitzRating:   1500,
		RapidRating:   1500,
		ClassicRating: 1500,
		Status:        "available",
	}

	res, err := InsertPlayerQuery(db, player)
	if err != nil {
		return nil, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	player.ID = lastInsertID

	return player, nil
}
