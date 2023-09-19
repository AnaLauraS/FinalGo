package handler

import (
	"errors"
	"finalgo/internal/odontologo"
	"finalgo/internal/turno"
	"finalgo/pkg/web"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// creo la estructura del controlador, inyectando el service
type odontologoHandler struct {
	s odontologo.Service
	turnoService turno.Service
}

// funcion para instanciar el controlador
func NewodOntologoHandler(s odontologo.Service) *odontologoHandler {
	return &odontologoHandler{
		s: s,
	}
}

// POST --> agregar odontologo
// Odontologo godoc
// @Summary Create Odontologo
// @Description Create a new odontologo
// @Tags odontologo
// @Accept json
// @Produce json
// @Success 201 {object} web.response
// @Failure 400 {object} web.errorResponse
// @Failure 500 {object} web.errorResponse
// @Router /odontologo [post]
func (h *odontologoHandler) CreateOdontologo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var odontologo odontologo.OdontologoRequest

		err := c.Bind(&odontologo)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// valido la existencia de datos clave
		valid, err := validateEmptys2(odontologo)
		if !valid {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		p, err := h.s.CreateOdontologo(c, odontologo)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}
		web.OkResponse(c, 201, p)
	}
}

// validateEmptys valida que los campos claves no esten vacios
func validateEmptys2(odontologo odontologo.OdontologoRequest) (bool, error) {
	if (odontologo.Apellido == "" || odontologo.Nombre == "" || odontologo.Matricula == "") {
		return false, errors.New("No se permiten los campos apellido, nombre y matricula vacíos")
	}
	return true, nil
}

// GET --> traer odontologo por id
// Odontologo godoc
// @Summary get odontologo
// @Description Get odontologo by id
// @Tags odontologo
// @Param id path int true "id del odontologo"
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Failure 400 {object} web.errorResponse
// @Failure 500 {object} web.errorResponse
// @Router /odontologo/:id [get]
func (h *odontologoHandler) GetOdontologoByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// obtengo el ID que pasaron por parámetro
		idQuery := ctx.Param("id")

		// valido el dato ingresado
		if idQuery != "" {

			// lo transformo a int y si me trae error, lo devuelvo
			id, err := strconv.Atoi(idQuery)
			if err != nil {
				if err != nil {
					web.ErrorResponse(ctx, http.StatusBadRequest)
					return
				}
			}

			// obtengo el odontologo
			odontologo, err := h.s.GetOdontologoByID(ctx, id)
			if err != nil {
				web.ErrorResponse(ctx, http.StatusNotFound)
				return
			}
			web.OkResponse(ctx, http.StatusOK, odontologo)
			return
		}

		// pero si no me pasaron id, devuelvo igual todos los odontologos
		odontologos, err := h.s.GetAll(ctx)
		if err != nil {
			web.ErrorResponse(ctx, http.StatusInternalServerError)
			return
		}
		web.OkResponse(ctx, http.StatusOK, odontologos)
	}
}

// PUT --> actualiza completo un odontologo
// Odontologo godoc
// @Summary update odontologo
// @Description Update odontologo by id
// @Tags odontologo
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Failure 400 {object} web.errorResponse
// @Failure 500 {object} web.errorResponse
// @Router /odontologo/:id [put]
func (h *odontologoHandler) UpdateOdontologo() gin.HandlerFunc {
	return func(c *gin.Context) {
		// valido id
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// verifico el json a enviar
		var odontologo odontologo.OdontologoRequest
		err = c.ShouldBindJSON(&odontologo)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// valido la existencia de datos clave
		valid, err := validateEmptys2(odontologo)
		if !valid {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// llamo al servicio para actualizar al odontologo
		o, err := h.s.UpdateOdontologo(c, odontologo, id)
		if err != nil {
			web.ErrorResponse(c, http.StatusInternalServerError)
			return
		}

		web.OkResponse(c, http.StatusOK, o)
	}
}

// PATCH --> actualiza parcial un odontologo
// Odontologo godoc
// @Summary update odontologo for field
// @Description Update odontologo for field
// @Tags odontologo
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Failure 400 {object} web.errorResponse
// @Failure 500 {object} web.errorResponse
// @Router /odontologo/patch/:id [patch]
func (h *odontologoHandler) UpdateOdontologoForField() gin.HandlerFunc {
	return func(c *gin.Context) {
		// valido id
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// le paso los query params
		apellidoQuery := c.Query("apellido")
		nombreQuery := c.Query("nombre")
		matriculaQuery := c.Query("matricula")

		// obtengo los datos del odontologo original
		odontologoOriginal, err := h.s.GetOdontologoByID(c, id)
		if err != nil {
			web.ErrorResponse(c, http.StatusNotFound)
			return
		}

		// creo el odontologo request con los datos del original
		odontologoRequest := odontologo.OdontologoRequest{
			Apellido:  odontologoOriginal.Apellido,
			Nombre:    odontologoOriginal.Nombre,
			Matricula: odontologoOriginal.Matricula,
		}

		// verifico si los campos tienen datos, los casteo y se los asigno al odontologo request
		if nombreQuery != "" {
			odontologoRequest.Nombre = nombreQuery
		}
		if apellidoQuery != "" {
			odontologoRequest.Apellido = apellidoQuery
		}
		if matriculaQuery != "" {
			odontologoRequest.Matricula = matriculaQuery
		}

		// llamo al metodo de actualizar odontologo, usando el odontologoRequest
		o, err := h.s.UpdateOdontologo(c, odontologoRequest, id)
		if err != nil {
			web.ErrorResponse(c, http.StatusInternalServerError)
			return
		}

		web.OkResponse(c, http.StatusOK, o)
	}
}

// DELETE --> elimina un odontologo
// Odontologo godoc
// @Summary delete odontologo
// @Description Delete odontologo by id
// @Tags odontologo
// @Param id path int true "id del odontologo"
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Failure 400 {object} web.errorResponse
// @Failure 500 {object} web.errorResponse
// @Router /odontologo/:id [delete]
func (h *odontologoHandler) DeleteOdontologo() gin.HandlerFunc {
	return func(c *gin.Context) {
		// valido id
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// busco los turnos asociados y se los elimino tambien
		turnos, errorT := h.turnoService.GetTurnoByOdontologo(c, id)
		if errorT == nil {
			for _, turno := range turnos {
				err := h.turnoService.DeleteTurno(c, turno.ID)
				if err != nil {
					web.ErrorResponse(c, http.StatusBadRequest)
				}
			}
		}
		// luego del borrado de todos los turnos, avanza con el borrado del odontologo

		// si falla el delete, es porque el ID era invalido
		err = h.s.DeleteOdontologo(c, id)
		if err != nil {
			web.ErrorResponse(c, http.StatusNotFound)
			return
		}
		respuesta := "Odontologo de ID " + c.Param("id") + " eliminado"
		web.OkResponse(c, http.StatusOK, respuesta)
	}
}
