package psql

import (
	"database/sql"
	"fmt"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"github.com/jllopis/try5/user"
	_ "github.com/lib/pq"
	"github.com/mgutz/dat/v1"
	"github.com/mgutz/dat/v1/sqlx-runner"
)

// PsqlStore hold the connection to the database and it has the properties
// to connect to it.
type PsqlStore struct {
	C      *runner.Connection
	Status int
	PsqlStoreOptions
}

// PsqlStoreOptions host the options for the databasef
type PsqlStoreOptions struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxIdleConns int
	MaxOpenConns int
}

// OpenPgSQLStore inicializa la conexión con la base de datos utilizando la configuración por defecto. Es una función variádica que acepta el paso de funciones del tipo func(*PsqlStore) error para la configuración
func OpenPgSQLStore(opts *PsqlStoreOptions) (*PsqlStore, error) {
	r := &PsqlStore{
		Status: 0,
	}
	r.PsqlStoreOptions = *opts
	if r.SSLMode == "" {
		r.SSLMode = "disable"
	}
	if r.MaxIdleConns == 0 {
		r.MaxIdleConns = 20
	}
	if r.MaxOpenConns == 0 {
		r.MaxIdleConns = 30
	}
	ds := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s port=%d", r.User, r.DBName, r.SSLMode, r.Password, r.Host, r.Port)
	db, err := sql.Open("postgres", ds)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(r.MaxIdleConns)
	db.SetMaxOpenConns(r.MaxOpenConns)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// set this to enable interpolation
	dat.EnableInterpolation = true
	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = false
	r.C = runner.NewConnection(db, "postgres")

	return r, nil
}

func (s *PsqlStore) LoadUser(uuid string) (*user.User, error) {
	var res *user.User
	if err := s.C.Select("*").From("users").Where("uuid=$1 AND deleted IS NULL", uuid).QueryStruct(&res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *PsqlStore) LoadAllUsers() ([]*user.User, error) {
	var res []*user.User
	if err := s.C.Select("*").From("users").Where("deleted IS NULL").QueryStructs(&res); err != nil {
		return nil, err
	}
	return res, nil
}

// SaveUser creates a new user if user.UID has zero value or updates the user otherwise.
func (s *PsqlStore) SaveUser(user *user.User) (*user.User, error) {
	now := time.Now().UTC()
	user.Updated = now
	switch user.UID {
	case "":
		user.UID = uuid.New()
		user.Created = now
		if err := s.C.InsertInto("users").Blacklist("id").Record(user).Returning("id").QueryScalar(&user.ID); err != nil {
			return user, err
		}
	default:
		if _, err := s.C.Update("users").SetBlacklist(&user, "id", "created").Where("id=$1", user.ID).Exec(); err != nil {
			return nil, err
		}
	}
	return user, nil
}

// DeleteUser elimina de la base de datos el User cuyo id coincide con id.
// Si la petición tiene éxito, devuelve el número de registros eliminados.
//
// Si aparece un error, devuelve el error del tipo *pq.Error
func (s *PsqlStore) DeleteUser(uuid string) (int, error) {
	var err error
	var res *dat.Result
	if res, err = s.C.DeleteFrom("users").Where("uuid = $1", uuid).Exec(); err != nil {
		return 0, err
	}
	if res.RowsAffected == 0 {
		return 0, nil
	}
	return int(res.RowsAffected), nil
}
