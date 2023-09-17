package turno

import "time"

// creamos la estructura del turno.
type Turno struct {
	ID           int       `json:"id"`
	IdOdontologo int       `json:"id_odontologo"`
	IdPaciente   int       `json:"id_paciente"`
	FechaHora    time.Time `json:"fecha_hora"`
	Descripcion  string    `json:"descripcion"`
}

// creamos la misma estructura de turno para las solicitudes por API.
type TurnoRequest struct {
	IdOdontologo int       `json:"id_odontologo"`
	IdPaciente   int       `json:"id_paciente"`
	FechaHora    time.Time `json:"fecha_hora"`
	Descripcion  string    `json:"descripcion"`
}

type TurnoDniMatriculaRequest struct {
	MatriculaOdontologo string    `json:"matricula_odontologo"`
	DniPaciente         string    `json:"dni_paciente"`
	FechaHora           time.Time `json:"fecha_hora"`
	Descripcion         string    `json:"descripcion"`
}
