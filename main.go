package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/pairing", pairingGameRequestHandler)
	http.HandleFunc("/cancel", cancelGameRequestHandler)
	http.HandleFunc("/game-ended", gameEndedHandler)
	http.HandleFunc("/create-game", createGameHandler)
	http.HandleFunc("/fetch-game", fetchGameByIDHandler)
	http.HandleFunc("/fetch-player", fetchPlayerByIDHandler)
	http.HandleFunc("/create-player", createNewPlayerHandler)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
