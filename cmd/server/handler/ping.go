package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type pingHandler struct {
}

func NewPingHandler() *pingHandler {
	return &pingHandler{}
}

// Ping godoc
// @Summary ping
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} pong
// @Router /ping [get]
func (c *pingHandler) Ping() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}

}

