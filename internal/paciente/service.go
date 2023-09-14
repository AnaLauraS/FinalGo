package paciente

import (
	"context"
	"log"
)

// defino la interfaz para que se apliquen siempre todos los métodos
type Service interface {
	GetPacienteByID(ctx context.Context, id int) (Paciente, error)
	GetAll(ctx context.Context) ([]Paciente, error)
	CreatePaciente(ctx context.Context, p PacienteRequest) (Paciente, error)
	UpdatePaciente(ctx context.Context, p PacienteRequest, id int) (Paciente, error)
	DeletePaciente(ctx context.Context, id int) error
}

// estrucutra service que contará con un repositorio
type service struct {
	r Repository
}

// función para instanciar service
func NewService(r Repository) Service {
	return &service{r}
}


// aplico todos los métodos de la interfaz, llamando a su correspondiente método del repositorio (todos conservan el mismo nombre):

func (s *service) GetAll(ctx context.Context) ([]Paciente, error) {
	pacientes, err := s.r.GetAll(ctx)
	if err != nil {
		log.Println("log de error en service de pacientes", err.Error())
		return []Paciente{}, ErrEmptyList
	}
	return pacientes, nil
}

func (s *service) GetPacienteByID(ctx context.Context, id int) (Paciente, error) {
	p, err := s.r.GetPacienteByID(ctx, id)
	if err != nil {
		log.Println("log de error por paciente inexistente", err.Error())
		return Paciente{}, ErrNotFound
	}
	return p, nil
}

func (s *service) CreatePaciente(ctx context.Context, pacienteRequest PacienteRequest) (Paciente, error) {
	// uso la estructura de request para mejor manejo de campos (no tiene el ID), llamando a una función que lo transforma en el dato que requiere la DB
	paciente := requestToPaciente(pacienteRequest)
	response, err := s.r.CreatePaciente(ctx, paciente)
	if err != nil {
		log.Println("error al crear paciente")
		return Paciente{}, ErrExec
	}
	return response, nil
}

func (s *service) DeletePaciente(ctx context.Context, id int) error {
	err := s.r.DeletePaciente(ctx, id)
	if err != nil {
		log.Println("log de error borrado de paciente", err.Error())
		return ErrNotFound
	}
	return nil
}

// este método está preparado para ser usado como PATCH o como PUT, se le deberá pasar desde el handler el paciente completo
func (s *service) UpdatePaciente(ctx context.Context, p PacienteRequest, id int) (Paciente, error) {
	// uso la estructura de request para mejor manejo de campos (no tiene el ID), llamando a una función que lo transforma en el dato que requiere la DB
	paciente := requestToPaciente(p)
	paciente.ID = id
	response, err := s.r.UpdatePaciente(ctx, paciente)
	if err != nil {
		log.Println("error al actualizar paciente")
		return Paciente{}, ErrExec
	}
	return response, nil
}


// función para transformar request en la estructura definida en GO
func requestToPaciente(pacienteRequest PacienteRequest) Paciente {
	var paciente Paciente
	paciente.Nombre = pacienteRequest.Nombre
	paciente.Apellido = pacienteRequest.Apellido
	paciente.Domicilio = pacienteRequest.Domicilio
	paciente.DNI = pacienteRequest.DNI
	paciente.Alta = pacienteRequest.Alta
	return paciente
}