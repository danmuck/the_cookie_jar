package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var WebSocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// This is to allow all origins, you might want to check it more carefully in production
		return true
	},
}

func RouteError(context *gin.Context, errorMessage string) {
	context.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
		"IsLoggedIn":   false,
		"ErrorMessage": errorMessage,
	})
}

func RouteIPLimit(context *gin.Context, waitTime int, unit string) {
	context.String(http.StatusTooManyRequests, "You have sent too many requests, please wait %d %s.", waitTime, unit)
}
