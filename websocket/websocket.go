package websocket

import (
	"fmt"
	"icc-backend-test/model"
)

// SendURLToPlayers sends the URL of the game to the players using WebSockets. When the players receive the URL, they can start playing the game. The frontend will use the URL to redirect the players to the game.
func SendURLToPlayers(player1, player2 *model.Player, url string) {
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
