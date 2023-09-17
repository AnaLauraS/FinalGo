package odontologo

import (
	"context"
	"log"
)

// defino la interfaz para que se apliquen siempre todos los métodos
type Service interface {
	GetOdontologoByID(ctx context.Context, id int) (Odontologo, error)
	GetAll(ctx context.Context) ([]Odontologo, error)
	CreateOdontologo(ctx context.Context, o OdontologoRequest) (Odontologo, error)
	UpdateOdontologo(ctx context.Context, o OdontologoRequest, id int) (Odontologo, error)
	DeleteOdontologo(ctx context.Context, id int) error
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

func (s *service) GetAll(ctx context.Context) ([]Odontologo, error) {
	odontologos, err := s.r.GetAll(ctx)
	if err != nil {
		log.Println("log de error en service de odontologo", err.Error())
		return []Odontologo{}, ErrEmptyList
	}
	return odontologos, nil
}

func (s *service) GetOdontologoByID(ctx context.Context, id int) (Odontologo, error) {
	o, err := s.r.GetOdontologoByID(ctx, id)
	if err != nil {
		log.Println("log de error por odontologo inexistente", err.Error())
		return Odontologo{}, ErrNotFound
	}
	return o, nil
}

func (s *service) CreateOdontologo(ctx context.Context, odontologoRequest OdontologoRequest) (Odontologo, error) {
	// uso la estructura de request para mejor manejo de campos (no tiene el ID), llamando a una función que lo transforma en el dato que requiere la DB
	odontologo := requestToOdontologo(odontologoRequest)
	response, err := s.r.CreateOdontologo(ctx, odontologo)
	if err != nil {
		log.Println("error al crear Odontologo")
		return Odontologo{}, ErrExec
	}
	return response, nil
}

func (s *service) DeleteOdontologo(ctx context.Context, id int) error {
	err := s.r.DeleteOdontologo(ctx, id)
	if err != nil {
		log.Println("log de error borrado de Odontologo", err.Error())
		return ErrNotFound
	}
	return nil
}

// este método está preparado para ser usado como PATCH o como PUT, se le deberá pasar desde el handler el Odontologo completo
func (s *service) UpdateOdontologo(ctx context.Context, odontologoRequest OdontologoRequest, id int) (Odontologo, error) {
	// uso la estructura de request para mejor manejo de campos (no tiene el ID), llamando a una función que lo transforma en el dato que requiere la DB
	odontologo := requestToOdontologo(odontologoRequest)
	odontologo.ID = id
	response, err := s.r.UpdateOdontologo(ctx, odontologo)
	if err != nil {
		log.Println("error al actualizar odontologo")
		return Odontologo{}, ErrExec
	}
	return response, nil
}

// función para transformar request en la estructura definida en GO
func requestToOdontologo(odontologoRequest OdontologoRequest) Odontologo {
	var odontologo Odontologo
	odontologo.Apellido = odontologoRequest.Apellido
	odontologo.Nombre = odontologoRequest.Nombre
	odontologo.Matricula = odontologoRequest.Matricula

	return odontologo
}
