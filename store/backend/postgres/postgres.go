package psql

import (
	"database/sql"
	"fmt"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"github.com/jllopis/try5/account"
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

var (
	notDeleted = dat.NewScope(
		"WHERE deleted IS NULL", nil)
)

func (s *PsqlStore) LoadAccount(uuid string) (*account.Account, error) {
	//var res *account.Account
	res := &account.Account{}
	if err := s.C.Select("*").From("accounts").Where("uid=$1 AND deleted IS NULL", uuid).QueryStruct(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *PsqlStore) LoadAllAccounts() ([]*account.Account, error) {
	var res []*account.Account
	if err := s.C.Select("*").From("accounts").ScopeMap(notDeleted, nil).QueryStructs(&res); err != nil {
		return nil, err
	}
	return res, nil
}

// Saveaccount creates a new account if account.UID has zero value or updates the account otherwise.
func (s *PsqlStore) SaveAccount(account *account.Account) (*account.Account, error) {
	now := time.Now().UTC()
	account.Updated = &now
	if account.Password != nil {
		account.UpdatePassword(*account.Password)
	}
	switch account.UID {
	case nil:
		u := uuid.New()
		account.UID = &u
		account.Created = &now
		if err := s.C.InsertInto("accounts").Blacklist("id").Record(account).Returning("id").QueryScalar(&account.ID); err != nil {
			return account, err
		}
	default:
		// TODO check if the provided uid effectively exists in the database prior to update
		if _, err := s.C.Update("accounts").SetBlacklist(account, "id", "uid", "created").Where("uid=$1", *account.UID).Exec(); err != nil {
			return nil, err
		}
	}
	return account, nil
}

// Deleteaccount elimina de la base de datos el account cuyo id coincide con id.
// Si la petición tiene éxito, devuelve el número de registros eliminados.
//
// Si aparece un error, devuelve el error del tipo *pq.Error
func (s *PsqlStore) DeleteAccount(uuid string) (int, error) {
	var err error
	var res *dat.Result
	if res, err = s.C.DeleteFrom("accounts").Where("uid = $1", uuid).Exec(); err != nil {
		return 0, err
	}
	if res.RowsAffected == 0 {
		return 0, nil
	}
	return int(res.RowsAffected), nil
}
