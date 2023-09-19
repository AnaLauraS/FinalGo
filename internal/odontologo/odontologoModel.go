package odontologo

// creamos la estructura de la entidad Odontologo. El ".json" especifica que deben serializarse y deserializarse al formato JSON
type Odontologo struct {
	ID        int    `json:"id"`
	Apellido  string `json:"apellido"`
	Nombre    string `json:"nombre"`
	Matricula string `json:"matricula"`
}

// creamos la misma estructura de Odontologo para las solicitudes por API o recibir datos de entrada.
type OdontologoRequest struct {
//	ID        int    `json:"id"`
	Apellido  string `json:"apellido"`
	Nombre    string `json:"nombre"`
	Matricula string `json:"matricula"`
}
