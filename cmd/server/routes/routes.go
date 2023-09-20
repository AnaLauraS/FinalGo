package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"finalgo/pkg/middleware"
	"finalgo/internal/odontologo"
	handler "finalgo/cmd/server/handler"
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
	r.buildPingRoutes()
}

// setGroup establece el grupo de enrutador.
func (r *router) setGroup() {
	r.routerGroup = r.engine.Group("/api/v1")
}

// buildOdontologoRoutes mapea todas las rutas para el dominio Odontologo.
func (r *router) buildOdontologoRoutes() {
	odontologoRepo := odontologo.NewRepositoryMySql(r.db)
	odontologoService := odontologo.NewService(odontologoRepo)
	turnoRepo := turno.NewRepositoryMySql(r.db)
	pacienteRepo := paciente.NewRepositoryMySql(r.db)
	pacienteService := paciente.NewService(pacienteRepo)
	turnoService := turno.NewService(turnoRepo, pacienteService, odontologoService)
	controladorOdontologo := handler.NewodOntologoHandler(odontologoService, turnoService)

	r.routerGroup.GET("/odontologos/:id", controladorOdontologo.GetOdontologoByID()) 
	r.routerGroup.POST("/odontologos", middleware.Authenticate(), controladorOdontologo.CreateOdontologo())
	r.routerGroup.PUT("/odontologos/:id", middleware.Authenticate(), controladorOdontologo.UpdateOdontologo())
	r.routerGroup.PATCH("/odontologos/:id", middleware.Authenticate(), controladorOdontologo.UpdateOdontologoForField())
	r.routerGroup.DELETE("/odontologos/:id", middleware.Authenticate(), controladorOdontologo.DeleteOdontologo())
}

// buildPacienteRoutes mapea todas las rutas para el dominio Paciente.
func (r *router) buildPacienteRoutes() {
	pacienteRepo := paciente.NewRepositoryMySql(r.db)
	pacienteService := paciente.NewService(pacienteRepo)
	odontologoRepo := odontologo.NewRepositoryMySql(r.db)
	odontologoService := odontologo.NewService(odontologoRepo)
	turnoRepo := turno.NewRepositoryMySql(r.db)
	turnoService := turno.NewService(turnoRepo, pacienteService, odontologoService)
	controladorPaciente := handler.NewPacienteHandler(pacienteService, turnoService)

	r.routerGroup.GET("/pacientes/:id", controladorPaciente.GetPacienteByID())
	r.routerGroup.POST("/pacientes", middleware.Authenticate(), controladorPaciente.CreatePaciente())
	r.routerGroup.PUT("/pacientes/:id", middleware.Authenticate(), controladorPaciente.UpdatePaciente())
	r.routerGroup.PATCH("/pacientes/:id", middleware.Authenticate(), controladorPaciente.UpdatePacienteForField())
	r.routerGroup.DELETE("/pacientes/:id", middleware.Authenticate(), controladorPaciente.DeletePaciente())
}

// buildTurnoRoutes mapea todas las rutas para el dominio Turno.
func (r *router) buildTurnoRoutes() {
	turnoRepo := turno.NewRepositoryMySql(r.db)
	pacienteRepo := paciente.NewRepositoryMySql(r.db)
	pacienteService := paciente.NewService(pacienteRepo)
	odontologoRepo := odontologo.NewRepositoryMySql(r.db)
	odontologoService := odontologo.NewService(odontologoRepo)
	turnoService := turno.NewService(turnoRepo, pacienteService, odontologoService)
	controladorTurno := handler.NewTurnoHandler(turnoService)

	r.routerGroup.GET("/turnos/:id", controladorTurno.GetTurnoByID())
	r.routerGroup.GET("/turnos/dni/:id", controladorTurno.GetTurnoByPaciente())
	r.routerGroup.POST("/turnos", middleware.Authenticate(), controladorTurno.CreateTurno())
	r.routerGroup.POST("/turnos/dni", middleware.Authenticate(), controladorTurno.CreateTurnoByDniAndMatricula())
	r.routerGroup.PUT("/turnos/:id", middleware.Authenticate(), controladorTurno.UpdateTurno())
	r.routerGroup.PATCH("/turnos/:id", middleware.Authenticate(), controladorTurno.UpdateTurnoForField())
	r.routerGroup.DELETE("/turnos/:id", middleware.Authenticate(), controladorTurno.DeleteTurno())
}

// API de prueba
func (r *router) buildPingRoutes() {
	r.routerGroup.GET("/ping", handler.NewPingHandler().Ping())
}