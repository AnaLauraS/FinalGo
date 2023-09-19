package odontologo

import (
	"context"
	"database/sql"
	"errors"
)

// Errores
var (
	ErrEmptyList = errors.New("la lista de Odontólogos esta vacia")
	ErrNotFound  = errors.New("Odontólogo no encontrado")
	ErrStatement = errors.New("sentencia incorrecta")
	ErrExec      = errors.New("ejecución SQL incorrecta")
	ErrLastId    = errors.New("error al obtener el último ID")
)

// Queries a usar en cada función
var (
	QueryInsert           = `INSERT INTO my_db.odontologo(apellido,nombre,matricula) VALUES(?,?,?,?,?)`
	QueryGetAll           = `SELECT id,apellido,nombre,matricula, alta FROM my_db.odontologo`
	QueryDelete           = `DELETE FROM my_db.odontologo WHERE id = ?`
	QueryGetById          = `SELECT id, apellido,nombre,matricula FROM my_db.odontologo WHERE id = ?`
	QueryUpdate           = `UPDATE my_db.odontologo SET apellido = ?,nombre = ?,matricula = ? WHERE id = ?`
	QueryGetIdByMatricula = `SELECT id FROM my_db.odontologo WHERE matricula = ?`
)

// defino la interfaz para que se apliquen siempre todos los métodos
type Repository interface {
	GetOdontologoByID(ctx context.Context, id int) (Odontologo, error)
	CreateOdontologo(ctx context.Context, o Odontologo) (Odontologo, error)
	UpdateOdontologo(ctx context.Context, o Odontologo) (Odontologo, error)
	GetAll(ctx context.Context) ([]Odontologo, error)
	DeleteOdontologo(ctx context.Context, id int) error
	GetOdontologoIdByMatricula(ctx context.Context, matricula string) (int, error)
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

// obtener todos los odontologos:
func (r *repository) GetAll(ctx context.Context) ([]Odontologo, error) {
	// ejecuto la query que trae todos los datos
	rows, err := r.db.Query(QueryGetAll)

	// si hay error de query, lo devuelvo
	if err != nil {
		return []Odontologo{}, ErrEmptyList
	}
	defer rows.Close()

	// creo el listado de odontologos
	var odontologos []Odontologo
	// voy poblando el listado de odontologos
	for rows.Next() {
		var odontologo Odontologo
		err := rows.Scan(
			&odontologo.ID,
			&odontologo.Apellido,
			&odontologo.Nombre,
			&odontologo.Matricula,
		)
		if err != nil {
			return []Odontologo{}, ErrExec
		}
		odontologos = append(odontologos, odontologo)
	}

	// verifico haber cargado bien todos los registros
	if err := rows.Err(); err != nil {
		return []Odontologo{}, ErrExec
	}

	// devuelvo el resultado
	return odontologos, nil
}

// obtener Odontologo por ID
func (r *repository) GetOdontologoByID(ctx context.Context, id int) (Odontologo, error) {
	// ejecuto la query de búsqueda por ID
	row := r.db.QueryRow(QueryGetById, id)

	// creo la variable que guarde (muestre) el resultado
	var odontologo Odontologo

	// verifico si obtengo algún error en los datos
	err := row.Scan(
		&odontologo.ID,
		&odontologo.Apellido,
		&odontologo.Nombre,
		&odontologo.Matricula,
	)

	// devuelvo el error o el odontologo
	if err != nil {
		return Odontologo{}, ErrNotFound
	}
	return odontologo, nil
}

// obtener Odontologo por ID
func (r *repository) GetOdontologoIdByMatricula(ctx context.Context, dni string) (int, error) {
	// ejecuto la query de búsqueda por ID
	row := r.db.QueryRow(QueryGetIdByMatricula, dni)

	// creo la variable que guarde (muestre) el resultado
	var odontologo Odontologo

	// verifico si obtengo algún error en los datos
	err := row.Scan(
		&odontologo.ID,
	)

	// devuelvo el error o el odontologo
	if err != nil {
		return Odontologo{}.ID, ErrNotFound
	}
	return odontologo.ID, nil
}

// crear Odontologo en BD
func (r *repository) CreateOdontologo(ctx context.Context, o Odontologo) (Odontologo, error) {
	// ejecuto la query
	statement, err := r.db.Prepare(QueryInsert)

	// verifico error de ejecución de query
	if err != nil {
		return Odontologo{}, ErrStatement
	}

	defer statement.Close()

	// paso los parámetros para que se ejecute la query
	result, err := statement.Exec(

		o.Apellido,
		o.Nombre,
		o.Matricula,
	)

	// verifico error de ejecución de query
	if err != nil {
		return Odontologo{}, ErrExec
	}

	// obtengo el ID del registro y lo devuelvo como dato
	lastId, err := result.LastInsertId()
	if err != nil {
		return Odontologo{}, ErrLastId
	}
	o.ID = int(lastId)
	return o, nil
}

// actualizar un registro
func (r *repository) UpdateOdontologo(ctx context.Context, o Odontologo) (Odontologo, error) {
	// preparo query para actualizar campos
	statement, err := r.db.Prepare(QueryUpdate)

	// por problemas de query, devuelve error
	if err != nil {
		return Odontologo{}, ErrStatement
	}
	defer statement.Close()

	// paso los parámetros para que se ejecute la query
	result, err := statement.Exec(
		o.Apellido,
		o.Nombre,
		o.Matricula,
	)

	// verifico error de parámetros
	if err != nil {
		return Odontologo{}, ErrStatement
	}

	// ejecuto query
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Odontologo{}, ErrExec
	}
	if rowsAffected < 1 {
		return Odontologo{}, ErrNotFound
	}

	return o, nil
}

// eliminar registro
func (r *repository) DeleteOdontologo(ctx context.Context, id int) error {
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
