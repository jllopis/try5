package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mgutz/dat/v1"
	"github.com/mgutz/dat/v1/sqlx-runner"
)

// StoreBackend hold the connection to the database and it has the properties
// to connect to it.
type StoreBackend struct {
	C        *runner.Connection
	Status   int
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// SetPgSQLHost realiza la configuración del host de la base de datos
func SetPgSQLHost(host string) func(*StoreBackend) error {
	return func(r *StoreBackend) error {
		r.Host = host
		return nil
	}
}

// SetPgSQLPort realiza la configuración del puerto en el host de la base de datos
func SetPgSQLPort(port int) func(*StoreBackend) error {
	return func(r *StoreBackend) error {
		r.Port = port
		return nil
	}
}

// SetPgSQLUser realiza la configuración del usuario de conexión a la base de datos
func SetPgSQLUser(user string) func(*StoreBackend) error {
	return func(r *StoreBackend) error {
		r.User = user
		return nil
	}
}

// SetPgSQLPassword realiza la configuración de la contraseña para la conexión a la base de datos
func SetPgSQLPassword(pass string) func(*StoreBackend) error {
	return func(r *StoreBackend) error {
		r.Password = pass
		return nil
	}
}

// SetPgSQLDBName realiza la configuración del nombre de la base de datos a la que conectaremos
func SetPgSQLDBName(dbname string) func(*StoreBackend) error {
	return func(r *StoreBackend) error {
		r.DBName = dbname
		return nil
	}
}

// SetPgSQLSSLMode realiza la configuración del modo SSL de la conexión a la base de datos
func SetPgSQLSSLMode(ssl string) func(*StoreBackend) error {
	return func(r *StoreBackend) error {
		r.SSLMode = ssl
		return nil
	}
}

// OpenPgSQLStore inicializa la conexión con la base de datos utilizando la configuración por defecto. Es una función variádica que acepta el paso de funciones del tipo func(*StoreBackend) error para la configuración
func OpenPgSQLStore(options ...func(*StoreBackend) error) (*StoreBackend, error) {
	r := &StoreBackend{
		Status:   0,
		Host:     "localhost",
		Port:     5432,
		User:     "",
		Password: "",
		DBName:   "",
		SSLMode:  "",
	}
	for _, option := range options {
		err := option(r)
		if err != nil {
			r = nil
			return nil, err
		}
	}
	if r.SSLMode == "" {
		r.SSLMode = "disable"
	}
	ds := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s port=%d", r.User, r.DBName, r.SSLMode, r.Password, r.Host, r.Port)
	db, err := sql.Open("postgres", ds)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(30)

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
