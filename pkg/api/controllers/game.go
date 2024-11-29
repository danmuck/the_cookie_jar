package controllers

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func GET_Game(c *gin.Context) {

}

// All open WebSockets, also a mutex to prevent race conditions
var openGamesSockets = make(map[*websocket.Conn]string)
var openGamesSocketsMutex sync.Mutex

// Broadcasts a message to all open sockets of a particular class ID
func broadcastGamesSockets(data []byte, id string) {
	openGamesSocketsMutex.Lock()
	defer openGamesSocketsMutex.Unlock()

	for socket := range openGamesSockets {
		// This open board isn't the same one we want to broadcast to
		if openGamesSockets[socket] != id {
			continue
		}

		err := socket.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			socket.Close()
			delete(openGamesSockets, socket)
		}
	}
}

func GET_GameWebSocket(c *gin.Context) {

}
