package main

import (
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"bitbucket.org/jllopis/getconf"
	"github.com/jllopis/try5/api"
	logger "github.com/jllopis/try5/log"
	_ "github.com/jllopis/try5/store/backend/boltdb"
	"github.com/jllopis/try5/store/manager"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/rs/cors"
)

// Config proporciona la configuración del servicio para ser utilizado por getconf
type Config struct {
	SslCert      string `getconf:"etcd app/try5/conf/sslcert, env TRY5_SSLCERT, flag sslcert"`
	SslKey       string `getconf:"etcd app/try5/conf/sslkey, env TRY5_SSLKEY, flag sslkey"`
	Port         string `getconf:"etcd app/try5/conf/port, env TRY5_PORT, flag port"`
	Origins      string `getconf:"etcd app/try5/conf/origins, env TRY5_ORIGINS, flag origins"`
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
	Revision    string
	config      *getconf.GetConf
	mainManager *manager.Manager
	//	apiCtx   *api.ApiContext
	verbose bool
)

func init() {
	//etcdURI := os.Getenv("TRY5_ETCD")
	//config = getconf.New(&Config{}, "TRY5", true, etcdURI)
	config = getconf.New(&Config{}, "TRY5", false, "")
	config.Parse()
	//dbPort := 5432
	//if p, err := config.GetInt("StorePort"); err == nil {
	//dbPort = int(p)
	//}
	//logger.LogD("Store", "port", dbPort)
	//	rs, err := psql.OpenPgSQLStore(&psql.PsqlStoreOptions{
	//		Host:     config.GetString("StoreHost"),
	//		Port:     dbPort,
	//		DBName:   config.GetString("StoreName"),
	//		User:     config.GetString("StoreUser",)
	//		Password: config.GetString("StorePass"),
	//	})

	/*
		apiCtx = &api.ApiContext{rs, r, securecookie.New(
			securecookie.GenerateRandomKey(64),
			securecookie.GenerateRandomKey(32))}
	*/
}

func main() {
	// Be sure we close the database when exit
	//	defer apiCtx.DB.Close()
	setupSignals()
	// Setup log
	debug := false
	if v, err := config.GetBool("Verbose"); err == nil && v {
		logger.SetLevel(5) // logrus.DebugLevel
		debug = true
		logger.LogD("set log level to DebugLevel")
	}

	setupStoreManager()

	// Setup api port
	port := config.GetString("Port")
	if port == "" {
		logger.LogW("can't get Port value from config", "USING:", 8000)
		port = "8000"
	}
	logger.LogI("Try5 API Server", "Version", Version, "Revision", Revision, "Build", BuildDate)
	logger.LogI("GetConf", "Version", getconf.Version())
	logger.LogI("Go", "Version", runtime.Version())
	logger.LogI("API Server", "Status", "started", "port", port)

	server := echo.New()
	server.SetDebug(debug)
	server.Use(mw.Recover())
	server.Use(mw.StripTrailingSlash())
	server.Use(mw.Logger())
	server.Get("/time", api.Time)
	// serve the V1 REST API from /api/v1
	apisrv := server.Group("/api/v1")
	// Gzip
	apisrv.Use(mw.Gzip())
	// Use CORS Handler and log every request to the api
	origins := strings.Split(config.GetString("Origins"), ",")
	if len(origins) == 0 || origins[0] == "" {
		origins = []string{"*"}
	}
	logger.LogD("setup cors", "allowed origins", origins)
	apisrv.Use(cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	setupAPIRoutes(apisrv)
	server.RunTLS(":"+port, config.GetString("SslCert"), config.GetString("SslKey"))
}

// setupAPIRoutes añade al router los puntos de acceso a los servicios ofrecidos
func setupAPIRoutes(apisrv *echo.Group) {
	// accounts
	apisrv.Get("/accounts", api.GetAllAccounts(mainManager))
	apisrv.Get("/accounts/:uid", api.GetAccountByID(mainManager))
	apisrv.Post("/accounts", api.NewAccount(mainManager))
	apisrv.Put("/accounts/:uid", api.UpdateAccount(mainManager))
	apisrv.Delete("/accounts/:uid", api.DeleteAccount(mainManager))
	// Keys
	apisrv.Get("/keys", api.GetAllKeys(mainManager))
	apisrv.Get("/keys/:kid", api.GetKey(mainManager))
	apisrv.Get("/accounts/:uid/keys", api.GetAccountKeys(mainManager))
	apisrv.Delete("/keys/:kid", api.DeleteKey(mainManager))
	// account jwt
	//apisrv.Get("/accounts/:uid/tokens", http.HandlerFunc(apiCtx.GetAccountTokens))

	// authentication
	apisrv.Post("/authenticate", api.Authenticate(mainManager))

	// JWT
	apisrv.Post("/jwt/token/:uid", http.HandlerFunc(apiCtx.NewJWTToken))
	apisrv.Post("/jwt/token/validate", http.HandlerFunc(apiCtx.ValidateToken))
	apisrv.Get("/jwt/token/validate", http.HandlerFunc(apiCtx.ValidateToken))
}

// setupSignals configura la captura de señales de sistema y actúa basándose en ellas
func setupSignals() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for sig := range sc {
			mainManager.Close()
			logger.LogI("signal.notify", "captured signal", sig, "stopping", true)
			os.Exit(1)
		}
	}()
}

func setupStoreManager() {
	timeout := 5 * time.Second
	if to, err := config.GetInt("StoreTimeout"); err == nil {
		logger.LogI("Store", "connect timeout (s)", to)
		timeout = time.Duration(to) * time.Second
	}
	var err error
	mainManager, err = manager.NewManager(map[string]interface{}{"path": config.GetString("StorePath"), "timeout": timeout}, "boltdb")
	if err != nil {
		logger.LogF("Could not create Manager: %s", err)
	}

	err = mainManager.Init()
	if err != nil {
		logger.LogF("Could not Init manager", "pkg", "main", "func", "Init()", "error", err)
	}
}
