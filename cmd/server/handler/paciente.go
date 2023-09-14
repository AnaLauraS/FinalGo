package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"finalgo/internal/paciente"
	"finalgo/pkg/web"
	"github.com/gin-gonic/gin"
)

// creo la estructura del controlador, inyectando el service
type pacienteHandler struct {
	s paciente.Service
}

// funcion para instanciar el controlador
func NewPacienteHandler(s paciente.Service) *pacienteHandler {
	return &pacienteHandler{
		s: s,
	}
}


// POST --> agregar paciente
func (h *pacienteHandler) CreatePaciente() gin.HandlerFunc {
	return func(c *gin.Context) {
		var paciente paciente.PacienteRequest

		err := c.ShouldBindJSON(&paciente)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// valido la existencia de datos clave
		valid, err := validateEmptys(paciente)
		if !valid {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		p, err := h.s.CreatePaciente(c, paciente)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}
		web.OkResponse(c, 201, p)
	}
}
// validateEmptys valida que los campos claves no esten vacios
func validateEmptys(paciente paciente.PacienteRequest) (bool, error) {
	if (paciente.Nombre == "" || paciente.Apellido == "" || paciente.DNI == "") {
		return false, errors.New("No se permiten los campos nombre, apellido y DNI vacíos")
	}
	return true, nil
}



// GET --> traer paciente por id
func (h *pacienteHandler) GetPacienteByID() gin.HandlerFunc {
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

			// obtengo el paciente
			paciente, err := h.s.GetPacienteByID(ctx, id)
			if err != nil {
				web.ErrorResponse(ctx, http.StatusNotFound)
				return
			}
			web.OkResponse(ctx, http.StatusOK, paciente)
			return
		}

		// pero si no me pasaron id, devuelvo igual todos los pacientes
		pacientes, err := h.s.GetAll(ctx)
		if err != nil {
			web.ErrorResponse(ctx, http.StatusInternalServerError)
			return
		}
		web.OkResponse(ctx, http.StatusOK, pacientes)
	}
}
	


// PUT --> actualiza completo un paciente
func (h *pacienteHandler) UpdatePaciente() gin.HandlerFunc {
	return func(c *gin.Context) {
		// valido id
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// verifico el json a enviar
		var paciente paciente.PacienteRequest
		err = c.ShouldBindJSON(&paciente)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// valido la existencia de datos clave
		valid, err := validateEmptys(paciente)
		if !valid {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// llamo al servicio para actualizar al paciente
		p, err := h.s.UpdatePaciente(c, paciente, id)
		if err != nil {
			web.ErrorResponse(c, http.StatusInternalServerError)
			return
		}

		web.OkResponse(c, http.StatusOK, p)
	}
}



// PATCH --> actualiza parcial un paciente
func (h *pacienteHandler) UpdatePacienteForField() gin.HandlerFunc {
	return func(c *gin.Context) {
		// valido id
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}

		// le paso los query params
		nombreQuery := c.Query("nombre")
		apellidoQuery := c.Query("apellido")
		domicilioQuery := c.Query("domiclio")
		dniQuery := c.Query("dni")
		altaQuery := c.Query("fecha_alta")

		// obtengo los datos del paciente original
		pacienteOriginal, err := h.s.GetPacienteByID(c, id)
		if err != nil {
			web.ErrorResponse(c, http.StatusNotFound)
			return
		}

		// creo el paciente request con los datos del original
		pacienteRequest := paciente.PacienteRequest{
			Nombre: pacienteOriginal.Nombre,
			Apellido: pacienteOriginal.Apellido,
			Domicilio: pacienteOriginal.Domicilio,
			DNI: pacienteOriginal.DNI,
			Alta: pacienteOriginal.Alta,
		}

		// verifico si los campos tienen datos, los casteo y se los asigno al paciente request
		if nombreQuery !="" {
			pacienteRequest.Nombre = nombreQuery
		}
		if apellidoQuery !="" {
			pacienteRequest.Apellido = apellidoQuery
		}
		if domicilioQuery !="" {
			pacienteRequest.Domicilio = domicilioQuery
		}
		if dniQuery !="" {
			pacienteRequest.DNI = dniQuery
		}
		if altaQuery !="" {
			fecha, err := time.Parse("2006-01-02", altaQuery)
			if err != nil {
				web.ErrorResponse(c, http.StatusBadRequest)
				return
			}
			pacienteRequest.Alta = fecha
		}
		
		// llamo al metodo de actualizar paciente, usando el pacienteRequest
		p, err := h.s.UpdatePaciente(c, pacienteRequest, id)
		if err != nil {
			web.ErrorResponse(c, http.StatusInternalServerError)
			return
		}

		web.OkResponse(c, http.StatusOK, p)
	}
}
	


// DELETE --> elimina un paciente
func (h *pacienteHandler) DeletePaciente() gin.HandlerFunc {
	return func(c *gin.Context) {
		// valido id
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.ErrorResponse(c, http.StatusBadRequest)
			return
		}
		
		// si falla el delete, es porque el ID era invalido
		err = h.s.DeletePaciente(c, id)
		if err != nil {
			web.ErrorResponse(c, http.StatusNotFound)
			return
		}
		respuesta := "Paciente de ID " + c.Param("id") + " eliminado"
		web.OkResponse(c, http.StatusOK, respuesta)
	}
}