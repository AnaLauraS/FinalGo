package paciente

import (
	"context"
	"database/sql"
	"errors"
)

// Errores
var (
	ErrEmptyList = errors.New("la lista de pacientes esta vacia")
	ErrNotFound  = errors.New("paciente no encontrado")
	ErrStatement = errors.New("sentencia incorrecta")
	ErrExec      = errors.New("ejecución SQL incorrecta")
	ErrLastId    = errors.New("error al obtener el último ID")
)

// Queries a usar en cada función
var (
	QueryInsert     = `INSERT INTO my_db.paciente(nombre, apellido, domicilio, dni, alta) VALUES(?,?,?,?,?)`
	QueryGetAll     = `SELECT id, nombre, apellido, domicilio, dni, alta FROM my_db.paciente`
	QueryDelete     = `DELETE FROM my_db.paciente WHERE id = ?`
	QueryGetById    = `SELECT id, nombre, apellido, domicilio, dni, alta FROM my_db.paciente WHERE id = ?`
	QueryUpdate     = `UPDATE my_db.paciente SET nombre = ?, apellido = ?, domicilio = ?, dni = ?, alta = ? WHERE id = ?`
	QueryGetIdByDni = `SELECT id FROM my_db.paciente WHERE dni = ?`
)

// defino la interfaz para que se apliquen siempre todos los métodos
type Repository interface {
	GetPacienteByID(ctx context.Context, id int) (Paciente, error)
	GetAll(ctx context.Context) ([]Paciente, error)
	CreatePaciente(ctx context.Context, p Paciente) (Paciente, error)
	UpdatePaciente(ctx context.Context, p Paciente) (Paciente, error)
	DeletePaciente(ctx context.Context, id int) error
	GetPacienteIDByDNI(ctx context.Context, dni string) (int, error)
}

// estructura repositorio con base de datos mysql
type repository struct {
	db *sql.DB
}

// NewRepositoryMySql instancia repositorio
func NewRepositoryMySql(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

// obtener todos los pacientes:
func (r *repository) GetAll(ctx context.Context) ([]Paciente, error) {
	// ejecuto la query que trae todos los datos
	rows, err := r.db.Query(QueryGetAll)

	// si hay error de query, lo devuelvo
	if err != nil {
		return []Paciente{}, ErrEmptyList
	}
	defer rows.Close()

	// creo el listado de pacientes
	var pacientes []Paciente
	// voy poblando el listado de pacientes
	for rows.Next() {
		var paciente Paciente
		err := rows.Scan(
			&paciente.ID,
			&paciente.Nombre,
			&paciente.Apellido,
			&paciente.Domicilio,
			&paciente.DNI,
			&paciente.Alta,
		)
		if err != nil {
			return []Paciente{}, ErrExec
		}
		pacientes = append(pacientes, paciente)
	}

	// verifico haber cargado bien todos los registros
	if err := rows.Err(); err != nil {
		return []Paciente{}, ErrExec
	}

	// devuelvo el resultado
	return pacientes, nil
}

// obtener pacientes por ID
func (r *repository) GetPacienteByID(ctx context.Context, id int) (Paciente, error) {
	// ejecuto la query de búsqueda por ID
	row := r.db.QueryRow(QueryGetById, id)

	// creo la variable que guarde (muestre) el resultado
	var paciente Paciente

	// verifico si obtengo algún error en los datos
	err := row.Scan(
		&paciente.ID,
		&paciente.Nombre,
		&paciente.Apellido,
		&paciente.Domicilio,
		&paciente.DNI,
		&paciente.Alta,
	)

	// devuelvo el error o el paciente
	if err != nil {
		return Paciente{}, ErrNotFound
	}
	return paciente, nil
}

// obtener ID del paciente por DNI
func (r *repository) GetPacienteIDByDNI(ctx context.Context, dni string) (int, error) {
	// ejecuto la query de búsqueda por ID
	row := r.db.QueryRow(QueryGetIdByDni, dni)

	// creo la variable que guarde (muestre) el resultado
	var paciente Paciente

	// verifico si obtengo algún error en los datos
	err := row.Scan(
		&paciente.ID,
	)

	// devuelvo el error o el paciente
	if err != nil {
		return Paciente{}.ID, ErrNotFound
	}
	return paciente.ID, nil
}

// crear paciente en BD
func (r *repository) CreatePaciente(ctx context.Context, paciente Paciente) (Paciente, error) {
	// ejecuto la query
	statement, err := r.db.Prepare(QueryInsert)

	// verifico error de ejecución de query
	if err != nil {
		return Paciente{}, ErrStatement
	}

	defer statement.Close()

	// paso los parámetros para que se ejecute la query
	result, err := statement.Exec(
		paciente.Nombre,
		paciente.Apellido,
		paciente.Domicilio,
		paciente.DNI,
		paciente.Alta,
	)

	// verifico error de ejecución de query
	if err != nil {
		return Paciente{}, ErrExec
	}

	// obtengo el ID del registro y lo devuelvo como dato
	lastId, err := result.LastInsertId()
	if err != nil {
		return Paciente{}, ErrLastId
	}
	paciente.ID = int(lastId)
	return paciente, nil
}

// actualizar un registro
func (r *repository) UpdatePaciente(ctx context.Context, paciente Paciente) (Paciente, error) {
	// preparo query para actualizar campos
	statement, err := r.db.Prepare(QueryUpdate)

	// por problemas de query, devuelve error
	if err != nil {
		return Paciente{}, ErrStatement
	}
	defer statement.Close()

	// paso los parámetros para que se ejecute la query
	result, err := statement.Exec(
		paciente.Nombre,
		paciente.Apellido,
		paciente.Domicilio,
		paciente.DNI,
		paciente.Alta,
		paciente.ID,
	)

	// verifico error de parámetros
	if err != nil {
		return Paciente{}, ErrStatement
	}

	// ejecuto query
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Paciente{}, ErrExec
	}
	if rowsAffected < 1 {
		return Paciente{}, ErrNotFound
	}

	return paciente, nil
}

// eliminar registro
func (r *repository) DeletePaciente(ctx context.Context, id int) error {
	// ejecuto query
	result, err := r.db.Exec(QueryDelete, id)

	// verifico error
	if err != nil {
		return ErrStatement
	}

	// verifico filas afectadas
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrExec
	}
	if rowsAffected < 1 {
		return ErrNotFound
	}

	return nil
}
