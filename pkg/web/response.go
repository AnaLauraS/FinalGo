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


// mensajes ante cada status como variable
var error400 = "error de datos enviados"
var error403 = "error de credenciales"
var error404 = "no encuentra elemento por error de datos enviados"
var error500 =  "problemas de servidor"
var errorDefault = "Internal Server Error"


// dos funciones: una de responseError (solo necesito pasarle el código de error porque el mensaje ya fue seteado) y otra de caso de exito
func ErrorResponse(c *gin.Context, status int){

	respuesta := errorResponse{
		Status: status,
	}

	// le defino los mensajes segun el status
	respuesta.Code = "Error"		
	switch status {
		case 400: respuesta.Message = error400
		case 403: respuesta.Message = error403
		case 404: respuesta.Message = error404
		case 500: respuesta.Message = error500
		default: respuesta.Message = error500
	}
	
	c.JSON(respuesta.Status, respuesta)
}

func OkResponse (c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}