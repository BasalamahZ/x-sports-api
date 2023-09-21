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
	"github.com/x-sports/internal/game"
	gamehttphandler "github.com/x-sports/internal/game/handler/http"
	gameservice "github.com/x-sports/internal/game/service"
	gamepgstore "github.com/x-sports/internal/game/store/postgresql"
	"github.com/x-sports/internal/match"
	matchhttphandler "github.com/x-sports/internal/match/handler/http"
	matchservice "github.com/x-sports/internal/match/service"
	matchpgstore "github.com/x-sports/internal/match/store/postgresql"
	"github.com/x-sports/internal/news"
	newshttphandler "github.com/x-sports/internal/news/handler/http"
	newsservice "github.com/x-sports/internal/news/service"
	newspgstore "github.com/x-sports/internal/news/store/postgresql"
	"github.com/x-sports/internal/team"
	teamhttphandler "github.com/x-sports/internal/team/handler/http"
	teamservice "github.com/x-sports/internal/team/service"
	teampgstore "github.com/x-sports/internal/team/store/postgresql"
	uploadhttphandler "github.com/x-sports/internal/upload/handler/http"
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

	// initialize game service
	var gameSvc game.Service
	{
		pgStore, err := gamepgstore.New(db)
		if err != nil {
			log.Printf("[game-api-http] failed to initialize game postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize game postgresql store: %s", err.Error())
		}

		gameSvc, err = gameservice.New(pgStore)
		if err != nil {
			log.Printf("[game-api-http] failed to initialize game service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize game service: %s", err.Error())
		}
	}

	// initialize team service
	var teamSvc team.Service
	{
		pgStore, err := teampgstore.New(db)
		if err != nil {
			log.Printf("[team-api-http] failed to initialize team postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize team postgresql store: %s", err.Error())
		}

		teamSvc, err = teamservice.New(pgStore)
		if err != nil {
			log.Printf("[team-api-http] failed to initialize team service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize team service: %s", err.Error())
		}
	}

	// initialize match service
	var matchSvc match.Service
	{
		pgStore, err := matchpgstore.New(db)
		if err != nil {
			log.Printf("[match-api-http] failed to initialize match postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize match postgresql store: %s", err.Error())
		}

		matchSvc, err = matchservice.New(pgStore)
		if err != nil {
			log.Printf("[match-api-http] failed to initialize match service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize match service: %s", err.Error())
		}
	}

	// initialize news service
	var newsSvc news.Service
	{
		pgStore, err := newspgstore.New(db)
		if err != nil {
			log.Printf("[news-api-http] failed to initialize news postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize news postgresql store: %s", err.Error())
		}

		newsSvc, err = newsservice.New(pgStore)
		if err != nil {
			log.Printf("[news-api-http] failed to initialize news service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize news service: %s", err.Error())
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

	// initialize game HTTP handler
	{
		identities := []gamehttphandler.HandlerIdentity{
			gamehttphandler.HandlerGames,
		}

		gameHTTP, err := gamehttphandler.New(gameSvc, adminSvc, identities)
		if err != nil {
			log.Printf("[game-api-http] failed to initialize game http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize game http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, gameHTTP)
	}

	// initialize team HTTP handler
	{
		identities := []teamhttphandler.HandlerIdentity{
			teamhttphandler.HandlerTeams,
		}

		teamHTTP, err := teamhttphandler.New(teamSvc, adminSvc, identities)
		if err != nil {
			log.Printf("[team-api-http] failed to initialize team http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize team http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, teamHTTP)
	}

	// initialize match HTTP handler
	{
		identities := []matchhttphandler.HandlerIdentity{
			matchhttphandler.HandlerMatch,
			matchhttphandler.HandlerMatchs,
		}

		matchHTTP, err := matchhttphandler.New(matchSvc, adminSvc, identities)
		if err != nil {
			log.Printf("[match-api-http] failed to initialize match http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize match http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, matchHTTP)
	}

	// initialize news HTTP handler
	{
		identities := []newshttphandler.HandlerIdentity{
			newshttphandler.HandlerNews,
			newshttphandler.HandlerNewss,
		}

		newsHTTP, err := newshttphandler.New(newsSvc, adminSvc, identities)
		if err != nil {
			log.Printf("[news-api-http] failed to initialize news http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize news http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, newsHTTP)
	}

	// initialize upload HTTP handler
	{
		identities := []uploadhttphandler.HandlerIdentity{
			uploadhttphandler.HandlerUpload,
		}

		uploadHTTP, err := uploadhttphandler.New(adminSvc, identities)
		if err != nil {
			log.Printf("[upload-api-http] failed to initialize upload http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize upload http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, uploadHTTP)
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
