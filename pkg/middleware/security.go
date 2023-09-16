package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	tokenHeaderName = "tokenPostman"
	envTokenKey     = "TOKEN"
	invalidUserMsg  = "Usuario inválido"
)

// Authenticate es un middleware que verifica si se proporciona el token correcto en la cabecera.
func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Obtener el token desde las variables de entorno.
		token := os.Getenv(envTokenKey)

		// Obtener el token de la cabecera de la solicitud.
		tokenHeader := ctx.GetHeader(tokenHeaderName)

		// Verificar si el token coincide.
		if tokenHeader != token {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": invalidUserMsg,
			})
			return
		}

		// Continuar con la solicitud si el token es válido.
		ctx.Next()
	}
}
