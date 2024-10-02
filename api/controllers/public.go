package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingPong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func Root(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":     "Index",
		"sub_title": "Future Homepage",
		"body":      "Some post text from a user that was in their recent post",
	})
}
