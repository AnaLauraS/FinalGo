package turno

import (
	"context"
	"database/sql"
	"errors"
)

// Errores
var (
	ErrEmptyList = errors.New("la lista de turnos esta vacia")
	ErrNotFound  = errors.New("turno no encontrado")
	ErrStatement = errors.New("sentencia incorrecta")
	ErrExec      = errors.New("ejecución SQL incorrecta")
	ErrLastId    = errors.New("error al obtener el último ID")
)

// Queries a usar en cada función
var (
	QueryInsert        = `INSERT INTO my_db.turnos(id_odontologo, id_paciente, fecha_hora, descripcion) VALUES(?,?,?,?)`
	QueryGetAll        = `SELECT id, id_odontologo, id_paciente, fecha_hora, descripcion FROM my_db.turnos`
	QueryDelete        = `DELETE FROM my_db.turnos WHERE id = ?`
	QueryGetById       = `SELECT id, id_odontologo, id_paciente, fecha_hora, descripcion FROM my_db.turnos WHERE id = ?`
	QueryUpdate        = `UPDATE my_db.turnos SET nombre = ?, apellido = ?, domicilio = ?, dni = ?, alta = ? WHERE id = ?`
	QueryGetByPaciente = `SELECT id, id_odontologo, id_paciente, fecha_hora, descripcion FROM my_db.turnos WHERE id_paciente = ?`
	QueryGetByOdontologo = `SELECT id, id_odontologo, id_paciente, fecha_hora, descripcion FROM my_db.turnos WHERE id_odontologo = ?`
)

// defino la interfaz para que se apliquen siempre todos los métodos
type Repository interface {
	GetTurnoByID(ctx context.Context, id int) (Turno, error)
	GetAll(ctx context.Context) ([]Turno, error)
	CreateTurno(ctx context.Context, p Turno) (Turno, error)
	UpdateTurno(ctx context.Context, p Turno) (Turno, error)
	DeleteTurno(ctx context.Context, id int) error
	GetTurnoByPaciente(ctx context.Context, id int) ([]Turno, error)
	GetTurnoByOdontologo(ctx context.Context, idOdontolog int) ([]Turno, error)
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

// obtener todos los turnos:
func (r *repository) GetAll(ctx context.Context) ([]Turno, error) {
	// ejecuto la query que trae todos los datos
	rows, err := r.db.Query(QueryGetAll)

	// si hay error de query, lo devuelvo
	if err != nil {
		return []Turno{}, ErrEmptyList
	}
	defer rows.Close()

	// creo el listado de turno
	var turnos []Turno
	// voy poblando el listado de turnos
	for rows.Next() {
		var turno Turno
		err := rows.Scan(
			&turno.ID,
			&turno.IdOdontologo,
			&turno.IdPaciente,
			&turno.FechaHora,
			&turno.Descripcion,
		)
		if err != nil {
			return []Turno{}, ErrExec
		}
		turnos = append(turnos, turno)
	}

	// verifico haber cargado bien todos los registros
	if err := rows.Err(); err != nil {
		return []Turno{}, ErrExec
	}

	// devuelvo el resultado
	return turnos, nil
}

// obtener turnos por ID
func (r *repository) GetTurnoByID(ctx context.Context, id int) (Turno, error) {
	// ejecuto la query de búsqueda por ID
	row := r.db.QueryRow(QueryGetById, id)

	// creo la variable que guarde (muestre) el resultado
	var turno Turno

	// verifico si obtengo algún error en los datos
	err := row.Scan(
		&turno.ID,
		&turno.IdOdontologo,
		&turno.IdPaciente,
		&turno.FechaHora,
		&turno.Descripcion,
	)

	// devuelvo el error o el turno
	if err != nil {
		return Turno{}, ErrNotFound
	}
	return turno, nil
}

// obtener turnos por ID del paciente
func (r *repository) GetTurnoByPaciente(ctx context.Context, id int) ([]Turno, error) {
	// ejecuto la query de búsqueda por ID
	row,err := r.db.Query(QueryGetByPaciente, id)

	// si hay error de query, lo devuelvo
	if err != nil {
		return []Turno{}, ErrEmptyList
	}
	defer row.Close()

	// creo la variable que guarde (muestre) el resultado
	var listadoTurno []Turno

	// verifico si obtengo algún error en los datos
	for row.Next() {
		var turno Turno
		err := row.Scan(
			&turno.ID,
			&turno.IdOdontologo,
			&turno.IdPaciente,
			&turno.FechaHora,
			&turno.Descripcion,
		)
		if err != nil {
			return []Turno{}, ErrExec
		}
		listadoTurno = append(listadoTurno, turno)
	}
	
	// verifico haber cargado bien todos los registros
	if err := row.Err(); err != nil {
		return []Turno{}, ErrExec
	}

	// devuelvo el resultado
	return listadoTurno, nil
}

// obtener turnos por ID del odontolog
func (r *repository) GetTurnoByOdontologo(ctx context.Context, id int) ([]Turno, error) {
	// ejecuto la query de búsqueda por ID
	row,err := r.db.Query(QueryGetByOdontologo, id)

	// si hay error de query, lo devuelvo
	if err != nil {
		return []Turno{}, ErrEmptyList
	}
	defer row.Close()

	// creo la variable que guarde (muestre) el resultado
	var listadoTurno []Turno

	// verifico si obtengo algún error en los datos
	for row.Next() {
		var turno Turno
		err := row.Scan(
			&turno.ID,
			&turno.IdOdontologo,
			&turno.IdPaciente,
			&turno.FechaHora,
			&turno.Descripcion,
		)
		if err != nil {
			return []Turno{}, ErrExec
		}
		listadoTurno = append(listadoTurno, turno)
	}
	
	// verifico haber cargado bien todos los registros
	if err := row.Err(); err != nil {
		return []Turno{}, ErrExec
	}

	// devuelvo el resultado
	return listadoTurno, nil
}

// crear turno en BD
func (r *repository) CreateTurno(ctx context.Context, turno Turno) (Turno, error) {
	// ejecuto la query
	statement, err := r.db.Prepare(QueryInsert)

	// verifico error de ejecución de query
	if err != nil {
		return Turno{}, ErrStatement
	}

	defer statement.Close()

	// paso los parámetros para que se ejecute la query
	result, err := statement.Exec(
		turno.IdOdontologo,
		turno.IdPaciente,
		turno.FechaHora,
		turno.Descripcion,
	)

	// verifico error de ejecución de query
	if err != nil {
		return Turno{}, ErrExec
	}

	// obtengo el ID del registro y lo devuelvo como dato
	lastId, err := result.LastInsertId()
	if err != nil {
		return Turno{}, ErrLastId
	}
	turno.ID = int(lastId)
	return turno, nil
}

// actualizar un registro
func (r *repository) UpdateTurno(ctx context.Context, turno Turno) (Turno, error) {
	// preparo query para actualizar campos
	statement, err := r.db.Prepare(QueryUpdate)

	// por problemas de query, devuelve error
	if err != nil {
		return Turno{}, ErrStatement
	}
	defer statement.Close()

	// paso los parámetros para que se ejecute la query
	result, err := statement.Exec(
		turno.IdOdontologo,
		turno.IdPaciente,
		turno.FechaHora,
		turno.Descripcion,
	)

	// verifico error de parámetros
	if err != nil {
		return Turno{}, ErrStatement
	}

	// ejecuto query
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Turno{}, ErrExec
	}
	if rowsAffected < 1 {
		return Turno{}, ErrNotFound
	}

	return turno, nil
}

// eliminar registro
func (r *repository) DeleteTurno(ctx context.Context, id int) error {
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
