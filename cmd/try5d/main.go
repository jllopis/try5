package main

import (
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"bitbucket.org/jllopis/getconf"
	"github.com/gorilla/securecookie"
	"github.com/jllopis/aloja"
	"github.com/jllopis/aloja/mw"
	"github.com/jllopis/try5/api"
	"github.com/jllopis/try5/store/backend/boltdb"
	"github.com/mgutz/logxi/v1"
	"github.com/rs/cors"
	"github.com/unrolled/render"
)

// Config proporciona la configuración del servicio para ser utilizado por getconf
type Config struct {
	SslCert      string `getconf:"etcd app/try5/conf/sslcert" env TRY5_SSLCERT, flag sslcert`
	SslKey       string `getconf:"etcd app/try5/conf/sslkey" env TRY5_SSLKEY, flag sslkey`
	Port         string `getconf:"etcd app/try5/conf/port, env TRY5_PORT, flag port"`
	Verbose      bool   `getconf:"etcd app/try5/conf/verbose, env TRY5_VERBOSE, flag verbose"`
	StorePath    string `getconf:"etcd app/try5/conf/storepath, env TRY5_STORE_PATH, flag storepath"`
	StoreTimeout int    `getconf:"etcd app/try5/conf/storetimeout, env TRY5_STORE_TIMEOUT, flag storetimeout"`
	//	StoreHost    string        `getconf:"etcd app/try5/conf/storehost, env TRY5_STORE_HOST, flag storehost"`
	//	StorePort    int           `getconf:"etcd app/try5/conf/storeport, env TRY5_STORE_PORT, flag storeport"`
	//	StoreName    string        `getconf:"etcd app/try5/conf/storename, env TRY5_STORE_NAME, flag storename"`
	//	StoreUser    string        `getconf:"etcd app/try5/conf/storeaccount, env TRY5_STORE_USER, flag storeuser"`
	//	StorePass    string        `getconf:"etcd app/try5/conf/storepass, env TRY5_STORE_PASS, flag storepass"`
}

var (
	// BuildDate holds the date the binary was built. It is valued at compile time
	BuildDate string
	// Version holds the version number of the build. It is valued at compile time
	Version string
	// Revision holds the git revision of the binary. It is valued at compile time
	Revision string
	config   *getconf.GetConf
	apiCtx   *api.ApiContext
	verbose  bool
	logger   log.Logger
)

func init() {
	//etcdURI := os.Getenv("TRY5_ETCD")
	//config = getconf.New(&Config{}, "TRY5", true, etcdURI)
	config = getconf.New(&Config{}, "TRY5", false, "")
	config.Parse()
	logger = log.New("try5api")
	//	dbPort := 5432
	//	if p, err := config.GetInt("StorePort"); err == nil {
	//		logger.Info("Store", "port", p)
	//		dbPort = int(p)
	//	}
	//	rs, err := psql.OpenPgSQLStore(&psql.PsqlStoreOptions{
	//		Host:     config.GetString("StoreHost"),
	//		Port:     dbPort,
	//		DBName:   config.GetString("StoreName"),
	//		User:     config.GetString("StoreUser"),
	//		Password: config.GetString("StorePass"),
	//	})
	timeout := 5 * time.Second
	if to, err := config.GetInt("StoreTimeout"); err == nil {
		logger.Info("Store", "connect timeout (s)", to)
		timeout = time.Duration(to) * time.Second
	}
	rs := bolt.NewBoltStore(&bolt.BoltStoreOptions{
		Dbpath:  config.GetString("StorePath"),
		Timeout: timeout,
	})
	if rs == nil {
		logger.Fatal("Cannot connect to file store", "db file path", rs.Dbpath)
	} else {
		logger.Info("Connected to store backend", "driver", "boltdb", "db file path", rs.Dbpath)
	}
	r := render.New(render.Options{
		Charset:    "UTF-8",
		PrefixXML:  []byte("<?xml version='1.0' encoding='UTF-8'?>"),
		IndentJSON: true,
	})
	apiCtx = &api.ApiContext{rs, r, securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32))}
}

func main() {
	// Be sure we close the database when exit
	defer apiCtx.DB.Close()
	setupSignals()
	port := config.GetString("Port")
	if port == "" {
		logger.Warn("can't get Port value from config", "USING:", 8000)
		port = "8000"
	}
	logger.Info("Try5 API Server", "Version", Version, "Revision", Revision, "Build", BuildDate)
	logger.Info("GetConf", "Version", getconf.Version())
	logger.Info("Go", "Version", runtime.Version())
	logger.Info("API Server", "Status", "started", "port", port)

	server := aloja.New().Port(port).SSLConf(config.GetString("SslCert"), config.GetString("SslKey"))
	// Use CORS Handler in every request and log every request
	cors := mw.CorsHandler(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"},
		AllowCredentials: true,
		Debug:            true,
	})
	server.AddGlobal(cors, mw.LogHandler)

	// serve the V1 REST API from /api/v1
	apisrv := server.NewSubrouter("/api/v1")
	setupAPIRoutes(apisrv)

	// run the server
	server.Run()
}

// setupAPIRoutes añade al router los puntos de acceso a los servicios ofrecidos
func setupAPIRoutes(apisrv *aloja.Subrouter) {
	// accounts
	apisrv.Get("/accounts", http.HandlerFunc(apiCtx.GetAllAccounts))
	apisrv.Get("/accounts/:uid", http.HandlerFunc(apiCtx.GetAccountByID))
	apisrv.Post("/accounts", http.HandlerFunc(apiCtx.NewAccount))
	apisrv.Put("/accounts/:uid", http.HandlerFunc(apiCtx.UpdateAccount))
	apisrv.Delete("/accounts/:uid", http.HandlerFunc(apiCtx.DeleteAccount))

	// authentication
	apisrv.Post("/authenticate", http.HandlerFunc(apiCtx.Authenticate))
}

// setupSignals configura la captura de señales de sistema y actúa basándose en ellas
func setupSignals() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for sig := range sc {
			logger.Info("signal.notify", "captured signal", sig, "stopping", true)
			os.Exit(1)
		}
	}()
}
