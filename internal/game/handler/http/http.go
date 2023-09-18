package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/x-sports/internal/admin"
	"github.com/x-sports/internal/game"
)

var (
	errUnknownConfig = errors.New("unknown config name")
)

// Handler contains admin HTTP-handlers.
type Handler struct {
	handlers map[string]*handler
	game     game.Service
	admin    admin.Service
}

// handler is the HTTP handler wrapper.
type handler struct {
	h        http.Handler
	identity HandlerIdentity
}

// HandlerIdentity denotes the identity of an HTTP hanlder.
type HandlerIdentity struct {
	Name string
	URL  string
}

// Followings are the known HTTP handler identities
var (
	// HandlerGames denotes HTTP handler to interact
	// with a games
	HandlerGames = HandlerIdentity{
		Name: "games",
		URL:  "/games",
	}
)

// New creates a new Handler.
func New(game game.Service, admin admin.Service, identities []HandlerIdentity) (*Handler, error) {
	h := &Handler{
		handlers: make(map[string]*handler),
		game:     game,
		admin:    admin,
	}

	// apply options
	for _, identity := range identities {
		if h.handlers == nil {
			h.handlers = map[string]*handler{}
		}

		h.handlers[identity.Name] = &handler{
			identity: identity,
		}

		handler, err := h.createHTTPHandler(identity.Name)
		if err != nil {
			return nil, err
		}

		h.handlers[identity.Name].h = handler
	}

	return h, nil
}

// createHTTPHandler creates a new HTTP handler that
// implements http.Handler.
func (h *Handler) createHTTPHandler(configName string) (http.Handler, error) {
	var httpHandler http.Handler
	switch configName {
	case HandlerGames.Name:
		httpHandler = &gamesHandler{
			game:  h.game,
			admin: h.admin,
		}
	default:
		return httpHandler, errUnknownConfig
	}
	return httpHandler, nil
}

// Start starts all HTTP handlers.
func (h *Handler) Start(multiplexer *mux.Router) error {
	for _, handler := range h.handlers {
		multiplexer.Handle(handler.identity.URL, handler.h)
	}
	return nil
}

// gameHTTP denotes user object in HTTP request or response
// body.
type gameHTTP struct {
	ID        *int64  `json:"id"`
	GameNames *string `json:"game_names"`
	GameIcons *string `json:"game_icons"`
}
