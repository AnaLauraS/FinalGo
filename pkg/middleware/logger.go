package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	basePath = "http://localhost:8080"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		path := ctx.Request.URL
		method := ctx.Request.Method

		ctx.Next()

		endTime := time.Now()
		elapsedTime := endTime.Sub(startTime)
		size := ctx.Writer.Size()

		log.Printf("Request:\n"+
			"  Method: %s\n"+
			"  Path: %s%s\n"+
			"  Time: %v\n"+
			"  Size: %d bytes",
			method, basePath, path, elapsedTime, size)
	}
}
