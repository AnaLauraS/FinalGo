package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"finalgo/pkg/middleware"
	"finalgo/cmd/server/handler"
	"finalgo/internal/odontologo"
	"finalgo/internal/paciente"
	"finalgo/internal/turno"
)

// Router es una interfaz que define los métodos que debe implementar cualquier enrutador.
type Router interface {
	MapRoutes()
}

// router es el enrutador Gin.
type router struct {
	engine      *gin.Engine
	routerGroup *gin.RouterGroup
	db          *sql.DB
}

// NewRouter crea un nuevo enrutador Gin.
func NewRouter(engine *gin.Engine, db *sql.DB) Router {
	return &router{
		engine: engine,
		db:     db,
	}
}

// MapRoutes mapea todas las rutas.
func (r *router) MapRoutes() {
	r.setGroup()
	r.buildOdontologoRoutes()
	r.buildPacienteRoutes()
	r.buildTurnoRoutes()
}

// setGroup establece el grupo de enrutador.
func (r *router) setGroup() {
	r.routerGroup = r.engine.Group("/api/v1")
}

// buildOdontologoRoutes mapea todas las rutas para el dominio Odontologo.
func (r *router) buildOdontologoRoutes() {
	odontologoRepo := odontologo.NewRepositoryMySql(r.db)
	odontologoService := odontologo.NewService(odontologoRepo)
	controladorOdontologo := odontologo.NewOdontologoHandler(odontologoService)

	r.routerGroup.GET("/odontologos", controladorOdontologo.GetOdontologos)
	r.routerGroup.GET("/odontologos/:id", controladorOdontologo.GetOdontologoByID)
	r.routerGroup.POST("/odontologos", controladorOdontologo.CreateOdontologo)
	r.routerGroup.PUT("/odontologos/:id", controladorOdontologo.UpdateOdontologo)
	r.routerGroup.DELETE("/odontologos/:id", controladorOdontologo.DeleteOdontologo)
}

// buildPacienteRoutes mapea todas las rutas para el dominio Paciente.
func (r *router) buildPacienteRoutes() {
	pacienteRepo := paciente.NewRepositoryMySql(r.db)
	pacienteService := paciente.NewService(pacienteRepo)
	controladorPaciente := paciente.NewPacienteHandler(pacienteService)

	r.routerGroup.GET("/pacientes", controladorPaciente.GetPacientes)
	r.routerGroup.GET("/pacientes/:id", controladorPaciente.GetPacienteByID)
	r.routerGroup.POST("/pacientes", controladorPaciente.CreatePaciente)
	r.routerGroup.PUT("/pacientes/:id", controladorPaciente.UpdatePaciente)
	r.routerGroup.DELETE("/pacientes/:id", controladorPaciente.DeletePaciente)
}

// buildTurnoRoutes mapea todas las rutas para el dominio Turno.
func (r *router) buildTurnoRoutes() {
	turnoRepo := turno.NewRepositoryMySql(r.db)
	turnoService := turno.NewService(turnoRepo, pacienteService, odontologoService)
	controladorTurno := turno.NewTurnoHandler(turnoService)

	r.routerGroup.GET("/turnos", controladorTurno.GetTurnos)
	r.routerGroup.GET("/turnos/:id", controladorTurno.GetTurnoByID)
	r.routerGroup.POST("/turnos", controladorTurno.CreateTurno)
	r.routerGroup.PUT("/turnos/:id", controladorTurno.UpdateTurno)
	r.routerGroup.DELETE("/turnos/:id", controladorTurno.DeleteTurno)
}
