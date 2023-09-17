package turno

import (
	"context"
	"finalgo/internal/odontologo"
	"finalgo/internal/paciente"
	"log"
)

// defino la interfaz para que se apliquen siempre todos los métodos
type Service interface {
	GetTurnoByID(ctx context.Context, id int) (Turno, error)
	GetAll(ctx context.Context) ([]Turno, error)
	CreateTurno(ctx context.Context, t TurnoRequest) (Turno, error)
	UpdateTurno(ctx context.Context, t TurnoRequest, id int) (Turno, error)
	DeleteTurno(ctx context.Context, id int) error
	GetTurnoByPaciente(ctx context.Context, dniPaciente string) (Turno, error)
	CreateTurnoByDniAndMatricula(ctx context.Context, t TurnoDniMatriculaRequest) (Turno, error)
}

// estrucutra service que contará con un repositorio
type service struct {
	r  Repository
	ps paciente.Service
	os odontologo.Service
}

// función para instanciar service
func NewService(r Repository, ps paciente.Service, os odontologo.Service) Service {
	return &service{
		r,
		ps,
		os,
	}
}

// aplico todos los métodos de la interfaz, llamando a su correspondiente método del repositorio (todos conservan el mismo nombre):

func (s *service) GetAll(ctx context.Context) ([]Turno, error) {
	turnos, err := s.r.GetAll(ctx)
	if err != nil {
		log.Println("log de error en service de turnos", err.Error())
		return []Turno{}, ErrEmptyList
	}
	return turnos, nil
}

func (s *service) GetTurnoByID(ctx context.Context, id int) (Turno, error) {
	p, err := s.r.GetTurnoByID(ctx, id)
	if err != nil {
		log.Println("log de error por turno inexistente", err.Error())
		return Turno{}, ErrNotFound
	}
	return p, nil
}

func (s *service) GetTurnoByPaciente(ctx context.Context, dniPaciente string) (Turno, error) {
	idPaciente, err := s.ps.GetPacienteIDByDNI(ctx, dniPaciente)
	if err != nil {
		log.Println("log de error por paciente inexistente", err.Error())
		return Turno{}, ErrNotFound
	}

	t, err := s.r.GetTurnoByPaciente(ctx, idPaciente)
	if err != nil {
		log.Println("log de error por turno inexistente", err.Error())
		return Turno{}, ErrNotFound
	}
	return t, nil
}

func (s *service) CreateTurno(ctx context.Context, turnoRequest TurnoRequest) (Turno, error) {
	// uso la estructura de request para mejor manejo de campos (no tiene el ID), llamando a una función que lo transforma en el dato que requiere la DB
	turno := requestToTurno(turnoRequest)
	response, err := s.r.CreateTurno(ctx, turno)
	if err != nil {
		log.Println("error al crear turno")
		return Turno{}, ErrExec
	}
	return response, nil
}

func (s *service) CreateTurnoByDniAndMatricula(ctx context.Context, t TurnoDniMatriculaRequest) (Turno, error) {
	// uso la estructura de request para mejor manejo de campos (no tiene el ID), llamando a una función que lo transforma en el dato que requiere la DB
	idPaciente, err := s.ps.GetPacienteIDByDNI(ctx, t.DniPaciente)
	if err != nil {
		log.Println("log de error por paciente inexistente", err.Error())
		return Turno{}, ErrNotFound
	}

	IdOdontologo, err := s.os.GetOdontologoIdByMatricula(ctx, t.MatriculaOdontologo)
	if err != nil {
		log.Println("log de error por odontologo inexistente", err.Error())
		return Turno{}, ErrNotFound
	}

	turnoRequest := TurnoRequest{
		IdOdontologo: IdOdontologo,
		IdPaciente:   idPaciente,
		FechaHora:    t.FechaHora,
		Descripcion:  t.Descripcion,
	}
	turno := requestToTurno(turnoRequest)
	response, err := s.r.CreateTurno(ctx, turno)
	if err != nil {
		log.Println("error al crear turno")
		return Turno{}, ErrExec
	}
	return response, nil
}

func (s *service) DeleteTurno(ctx context.Context, id int) error {
	err := s.r.DeleteTurno(ctx, id)
	if err != nil {
		log.Println("log de error borrado de turno", err.Error())
		return ErrNotFound
	}
	return nil
}

// este método está preparado para ser usado como PATCH o como PUT, se le deberá pasar desde el handler el turno completo
func (s *service) UpdateTurno(ctx context.Context, p TurnoRequest, id int) (Turno, error) {
	// uso la estructura de request para mejor manejo de campos (no tiene el ID), llamando a una función que lo transforma en el dato que requiere la DB
	turno := requestToTurno(p)
	turno.ID = id
	response, err := s.r.UpdateTurno(ctx, turno)
	if err != nil {
		log.Println("error al actualizar turno")
		return Turno{}, ErrExec
	}
	return response, nil
}

// función para transformar request en la estructura definida en GO
func requestToTurno(turnoRequest TurnoRequest) Turno {
	var turno Turno
	turno.IdOdontologo = turnoRequest.IdOdontologo
	turno.IdPaciente = turnoRequest.IdPaciente
	turno.FechaHora = turnoRequest.FechaHora
	turno.Descripcion = turnoRequest.Descripcion
	return turno
}
