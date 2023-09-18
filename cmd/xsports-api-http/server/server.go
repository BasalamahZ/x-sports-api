package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/x-sports/cmd/xsports-api-http/config"
	"github.com/x-sports/internal/admin"
	adminhttphandler "github.com/x-sports/internal/admin/handler/http"
	adminservice "github.com/x-sports/internal/admin/service"
	adminpgstore "github.com/x-sports/internal/admin/store/postgresql"
)

// Following constants are the possible exit code returned
// when running a server.
const (
	CodeSuccess = iota
	CodeBadConfig
	CodeFailServeHTTP
)

// Run creates a server and starts the server.
//
// Run returns a status code suitable for os.Exit() argument.
func Run() int {
	s, err := new()
	if err != nil {
		return CodeBadConfig
	}

	return s.start()
}

// server is the long-runnning application.
type server struct {
	srv      *http.Server
	handlers []handler
}

// handler provides mechanism to start HTTP handler. All HTTP
// handlers must implements this interface.
type handler interface {
	Start(multiplexer *mux.Router) error
}

// new creates and returns a new server.
func new() (*server, error) {
	s := &server{
		srv: &http.Server{
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	// connect to dabatabase
	db, err := sqlx.Connect("postgres", config.BaseConfig())
	if err != nil {
		log.Printf("[xsports-api-http] failed to connect database: %s\n", err.Error())
		return nil, fmt.Errorf("failed to connect database: %s", err.Error())
	}

	// initialize admin service
	var adminSvc admin.Service
	{
		pgStore, err := adminpgstore.New(db)
		if err != nil {
			log.Printf("[admin-api-http] failed to initialize admin postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize admin postgresql store: %s", err.Error())
		}

		svcOptions := []adminservice.Option{}
		svcOptions = append(svcOptions, adminservice.WithConfig(adminservice.Config{
			PasswordSalt:   os.Getenv("PasswordSalt"),
			TokenSecretKey: os.Getenv("TokenSecretKey"),
		}))

		adminSvc, err = adminservice.New(pgStore, svcOptions...)
		if err != nil {
			log.Printf("[tenant-api-http] failed to initialize admin service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize admin service: %s", err.Error())
		}
	}

	// initialize admin HTTP handler
	{
		identities := []adminhttphandler.HandlerIdentity{
			adminhttphandler.HandlerLogin,
		}

		adminHTTP, err := adminhttphandler.New(adminSvc, identities)
		if err != nil {
			log.Printf("[admin-api-http] failed to initialize admin http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize admin http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, adminHTTP)
	}

	return s, nil
}

// start starts the given server.
func (s *server) start() int {
	log.Println("[xsports-api-http] starting server...")

	// create multiplexer object
	rootMux := mux.NewRouter()
	appMux := rootMux.PathPrefix("/api/v1").Subrouter()

	// starts handlers
	for _, h := range s.handlers {
		if err := h.Start(appMux); err != nil {
			log.Printf("[xsports-api-http] failed to start handler: %s\n", err.Error())
			return CodeFailServeHTTP
		}
	}

	// endpoint checker
	appMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world! Auto Deploy On, INPO 60 RIBUNYA CAIRINN DONGG!!! @xsports")
	})

	// use middlewares to app mux only
	appMux.Use(corsMiddleware)

	// listen and serve
	log.Printf("[xsports-api-http] Server is running at %s:%s", os.Getenv("ADDRESS"), os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("ADDRESS"), os.Getenv("PORT")), rootMux))

	return CodeSuccess
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}