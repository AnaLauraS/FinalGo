package paciente

import "time"

// creamos la estructura de paciente. DNI tiene formato string porque no es un dato con el se deba hacer operaciones numéricas.
type Paciente struct {
	ID int `json:"id"`
	Nombre string `json:"nombre"`
	Apellido string `json:"apellido"`
	Domicilio string `json:"domicilio"`
	DNI string `json:"dni"`
	Alta time.Time `json:"fecha_alta"`
}

// creamos la misma estructura de paciente para las solicitudes por API.
type PacienteRequest struct {
	Nombre string `json:"nombre"`
	Apellido string `json:"apellido"`
	Domicilio string `json:"domicilio"`
	DNI string `json:"dni"`
	Alta time.Time `json:"fecha_alta"`
}