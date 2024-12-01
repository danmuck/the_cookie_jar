package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var WebSocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // This is to allow all origins, you might want to check it more carefully in production
	},
}

func RouteError(context *gin.Context, errorMessage string) {
	context.HTML(http.StatusBadRequest, "error.html", gin.H{
		"IsLoggedIn":   false,
		"ErrorMessage": errorMessage,
	})
}
