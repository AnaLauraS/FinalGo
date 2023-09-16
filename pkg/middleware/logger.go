package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	basePath = "http://localhost:8080"
)

// logger del middleware
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context){
		path := ctx.Request.URL
		time := time.Now()
		method := ctx.Request.Method
		var size int

		ctx.Next()

		if ctx.Writer != nil {
			size = ctx.Writer.Size()
		}

		fmt.Printf("Path: %s%s\nMethod: %s\nTime: %v\nSize: %d",basePath,path,method,time,size)
	}
}
