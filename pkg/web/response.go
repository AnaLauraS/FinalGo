package web

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type response struct {
	Data interface{} `json:"data"`
}

// dos funciones: una de responseError (solo necesito pasarle el código de error porque el mensaje ya fue seteado) y otra de caso de exito
func ErrorResponse(c *gin.Context, status int){

	respuesta := errorResponse{
		Status: status,
	}

	// le defino los mensajes segun el status
	respuesta.Code = "Error"		
	switch status {
		case 400: respuesta.Message = "error de datos enviados"
		case 403: respuesta.Message = "error de credenciales"
		case 404: respuesta.Message = "no encuentra elemento por error de datos enviados"
		case 500: respuesta.Message = "problemas de servidor"
		default: respuesta.Message = "Internal Server Error"
	}
	
	c.JSON(respuesta.Status, respuesta)
}

func OkResponse (c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}