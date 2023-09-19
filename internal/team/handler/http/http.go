package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/x-sports/internal/admin"
	"github.com/x-sports/internal/team"
)

var (
	errUnknownConfig = errors.New("unknown config name")
)

// Handler contains admin HTTP-handlers.
type Handler struct {
	handlers map[string]*handler
	team     team.Service
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
	// HandlerTeams denotes HTTP handler to interact
	// with a teams
	HandlerTeams = HandlerIdentity{
		Name: "teams",
		URL:  "/teams",
	}
)

// New creates a new Handler.
func New(team team.Service, admin admin.Service, identities []HandlerIdentity) (*Handler, error) {
	h := &Handler{
		handlers: make(map[string]*handler),
		team:     team,
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
	case HandlerTeams.Name:
		httpHandler = &teamsHandler{
			team:  h.team,
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

// teamHTTP denotes user object in HTTP request or response
// body.
type teamHTTP struct {
	ID        *int64  `json:"id"`
	TeamNames *string `json:"team_names"`
	GameID    *int64  `json:"game_id"`
	GameNames *string `json:"game_names"`
}
