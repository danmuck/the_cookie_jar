package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingPong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func Root(c *gin.Context) {
	c.String(http.StatusOK, "root controller")
}
