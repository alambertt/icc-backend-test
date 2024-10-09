package websocket

import (
	"icc-backend-test/model"
)

// SendURLToPlayers sends the URL of the game to the players using WebSockets. When the players receive the URL, they can start playing the game. The frontend will use the URL to redirect the players to the game.
func SendURLToPlayer(player *model.Player, url string) {
	// message := map[string]string{
	// 	"type": "game_url",
	// 	"url":  url,
	// }

	// Send URL to white player
	// if err := player.Conn.WriteJSON(message); err != nil {
	// 	fmt.Println("Error sending URL to player 1:", err)
	// }
}
