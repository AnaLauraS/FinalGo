package handler

import (
	"errors"
	"finalgo/internal/turno"
	"finalgo/pkg/web"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// creo la estructura del controlador, inyectando el service
type turnoHandler struct {
	s turno.Service
}

// funcion para instanciar el controlador
func NewTurnoHandler(s turno.Service) *turnoHandler {
	return &turnoHandler{
		s: s,
	}
}

// POST --> agregar turno
// Turno godoc
// @Summary Create Turno
// @Description Create a new turno
// @Tags turno
// @Accept json
// @Produce json
// @Success 201 {object} web.response
// @Failure 400 {object} web.errorResponse
// @Failure 500 {object} web.errorResponse
// @Router /turno [post]
func (h *turnoHandler) CreateTurno() gin.HandlerFunc {
	return func(c *gin.Context) {
		var turno turno.TurnoRequest

		err := c.ShouldBindJSON(&turno)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// valido la existencia de datos clave
		valid, err := validateTurnoEmptys(turno)
		if !valid {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		p, err := h.s.CreateTurno(c, turno)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}
		web.OkResponse(c, 201, p)
	}
}

// validateTurnoEmptys valida que los campos claves no esten vacios
func validateTurnoEmptys(turno turno.TurnoRequest) (bool, error) {
	if turno.IdOdontologo < 1 || turno.IdPaciente < 1 || turno.FechaHora.IsZero() {
		return false, errors.New("No se permiten los campos odontologo, paciente, fecha y hora vacíos")
	}
	return true, nil
}

// POST --> agregar turno con dni de paciente y matricula del odontologo
// Turno godoc
// @Summary Create Turno By DNI and Matricula
// @Description Create a new turno by DNI and Matricula
// @Tags turno
// @Accept json
// @Produce json
// @Success 201 {object} web.response
// @Failure 400 {object} web.errorResponse
// @Failure 500 {object} web.errorResponse
// @Router /turno/dni [post]
func (h *turnoHandler) CreateTurnoByDniAndMatricula() gin.HandlerFunc {
	return func(c *gin.Context) {
		var turno turno.TurnoDniMatriculaRequest

		err := c.ShouldBindJSON(&turno)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// valido la existencia de datos clave
		valid, err := validateTurnoEmptys2(turno)
		if !valid {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		t, err := h.s.CreateTurnoByDniAndMatricula(c, turno)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}
		web.OkResponse(c, 201, t)
	}
}

// validateTurnoEmptys valida que los campos claves no esten vacios
func validateTurnoEmptys2(turno turno.TurnoDniMatriculaRequest) (bool, error) {
	if turno.DniPaciente == "" || turno.MatriculaOdontologo == "" || turno.FechaHora.IsZero() {
		return false, errors.New("No se permiten los campos odontologo, paciente, fecha y hora vacíos")
	}
	return true, nil
}

// GET --> traer turno por id
// Turno godoc
// @Summary turno example
// @Description Get turno by id
// @Tags turno
// @Param id path int true "id del turno"
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Failure 400 {object} web.errorResponse
// @Failure 500 {object} web.errorResponse
// @Router /turno/:id [get]
func (h *turnoHandler) GetTurnoByID() gin.HandlerFunc {
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

			// obtengo el turno
			turno, err := h.s.GetTurnoByID(ctx, id)
			if err != nil {
				web.ErrorResponse(ctx, http.StatusNotFound)
				return
			}
			web.OkResponse(ctx, http.StatusOK, turno)
			return
		}

		// pero si no me pasaron id, devuelvo igual todos los turnos
		turnos, err := h.s.GetAll(ctx)
		if err != nil {
			web.ErrorResponse(ctx, http.StatusInternalServerError)
			return
		}
		web.OkResponse(ctx, http.StatusOK, turnos)
	}
}

// GET --> traer turno por dni del paciente
// Turno godoc
// @Summary get turno by dni
// @Description Get turno by dni
// @Tags turno
// @Param dni path string true "dni del paciente"
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Failure 400 {object} web.errorResponse
// @Failure 500 {object} web.errorResponse
// @Router /turno/dni/:id [get]
func (h *turnoHandler) GetTurnoByPaciente() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// obtengo el DNI que pasaron por parámetro
		dniQuery := ctx.Param("dni")

		// valido el dato ingresado
		if dniQuery != "" {
			// obtengo el turno
			turno, err := h.s.GetTurnoByPaciente(ctx, dniQuery)
			if err != nil {
				web.ErrorResponse(ctx, http.StatusNotFound)
				return
			}
			web.OkResponse(ctx, http.StatusOK, turno)
			return
		}

		// pero si no me pasaron id, devuelvo igual todos los turnos
		turnos, err := h.s.GetAll(ctx)
		if err != nil {
			web.ErrorResponse(ctx, http.StatusInternalServerError)
			return
		}
		web.OkResponse(ctx, http.StatusOK, turnos)
	}
}

// PUT --> actualiza completo un turno
// Turno godoc
// @Summary update turno
// @Description Update turno by id
// @Tags turno
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Failure 400 {object} web.errorResponse
// @Failure 500 {object} web.errorResponse
// @Router /turno/:id [put]
func (h *turnoHandler) UpdateTurno() gin.HandlerFunc {
	return func(c *gin.Context) {
		// valido id
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// verifico el json a enviar
		var turno turno.TurnoRequest
		err = c.ShouldBindJSON(&turno)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// valido la existencia de datos clave
		valid, err := validateTurnoEmptys(turno)
		if !valid {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// llamo al servicio para actualizar al turno
		p, err := h.s.UpdateTurno(c, turno, id)
		if err != nil {
			web.ErrorResponse(c, http.StatusInternalServerError)
			return
		}

		web.OkResponse(c, http.StatusOK, p)
	}
}

// PATCH --> actualiza parcial un turno
// Turno godoc
// @Summary update turno for field
// @Description Update turno for field
// @Tags turno
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Failure 400 {object} web.errorResponse
// @Failure 500 {object} web.errorResponse
// @Router /turno/patch/:id [patch]
func (h *turnoHandler) UpdateTurnoForField() gin.HandlerFunc {
	return func(c *gin.Context) {
		// valido id
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// le paso los query params
		odontologoQuery := c.Query("id_odontologo")
		pacienteQuery := c.Query("id_paciente")
		fechaHoraQuery := c.Query("fecha_hora")
		descripcionQuery := c.Query("descripcion")

		// obtengo los datos del turno original
		turnoOriginal, err := h.s.GetTurnoByID(c, id)
		if err != nil {
			web.ErrorResponse(c, http.StatusNotFound)
			return
		}

		// creo el turno request con los datos del original
		turnoRequest := turno.TurnoRequest{
			IdOdontologo: turnoOriginal.IdOdontologo,
			IdPaciente:   turnoOriginal.IdPaciente,
			FechaHora:    turnoOriginal.FechaHora,
			Descripcion:  turnoOriginal.Descripcion,
		}

		// verifico si los campos tienen datos, los casteo y se los asigno al turno request
		if odontologoQuery != "" {
			odontologoID, err := strconv.Atoi(odontologoQuery)
			if err != nil {
				web.ErrorResponse(c, http.StatusBadRequest)
				return
			}
			turnoRequest.IdOdontologo = odontologoID
		}
		if pacienteQuery != "" {
			pacienteID, err := strconv.Atoi(pacienteQuery)
			if err != nil {
				web.ErrorResponse(c, http.StatusBadRequest)
				return
			}
			turnoRequest.IdPaciente = pacienteID
		}
		if fechaHoraQuery != "" {
			fecha, err := time.Parse("2006-01-02 00:00", fechaHoraQuery)
			if err != nil {
				web.ErrorResponse(c, http.StatusBadRequest)
				return
			}
			turnoRequest.FechaHora = fecha
		}
		if descripcionQuery != "" {
			turnoRequest.Descripcion = descripcionQuery
		}

		// llamo al metodo de actualizar turno, usando el turnoRequest
		p, err := h.s.UpdateTurno(c, turnoRequest, id)
		if err != nil {
			web.ErrorResponse(c, http.StatusInternalServerError)
			return
		}

		web.OkResponse(c, http.StatusOK, p)
	}
}

// DELETE --> elimina un turno
// Turno godoc
// @Summary delete turno
// @Description Delete turno by id
// @Tags turno
// @Param id path int true "id del turno"
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Failure 400 {object} web.errorResponse
// @Failure 500 {object} web.errorResponse
// @Router /turno/:id [delete]
func (h *turnoHandler) DeleteTurno() gin.HandlerFunc {
	return func(c *gin.Context) {
		// valido id
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// si falla el delete, es porque el ID era invalido
		err = h.s.DeleteTurno(c, id)
		if err != nil {
			web.ErrorResponse(c, http.StatusNotFound)
			return
		}
		respuesta := "Turno de ID " + c.Param("id") + " eliminado"
		web.OkResponse(c, http.StatusOK, respuesta)
	}
}
