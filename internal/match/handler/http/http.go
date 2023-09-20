package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/x-sports/internal/admin"
	"github.com/x-sports/internal/match"
)

var (
	errUnknownConfig = errors.New("unknown config name")
)

// dateFormat denotes the standard date format used in
// match HTTP request and response.
var dateFormat = "2 January 2006"

// Handler contains admin HTTP-handlers.
type Handler struct {
	handlers map[string]*handler
	match    match.Service
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
	// HandlerMatch denotes HTTP handler to interact
	// with a match
	HandlerMatch = HandlerIdentity{
		Name: "match",
		URL:  "/matchs/{id}",
	}

	// HandlerMatchs denotes HTTP handler to interact
	// with a matchs
	HandlerMatchs = HandlerIdentity{
		Name: "matchs",
		URL:  "/matchs",
	}
)

// New creates a new Handler.
func New(match match.Service, admin admin.Service, identities []HandlerIdentity) (*Handler, error) {
	h := &Handler{
		handlers: make(map[string]*handler),
		match:    match,
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
	case HandlerMatch.Name:
		httpHandler = &matchHandler{
			match: h.match,
			admin: h.admin,
		}
	case HandlerMatchs.Name:
		httpHandler = &matchsHandler{
			match: h.match,
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

// matchHTTP denotes user object in HTTP request or response
// body.
type matchHTTP struct {
	ID              *int64   `json:"id"`
	TournamentNames *string  `json:"tournament_names"`
	GameID          *int64   `json:"game_id"`
	GameNames       *string  `json:"game_names"`
	TeamAID         *int64   `json:"team_a_id"`
	TeamANames      *string  `json:"team_a_names"`
	TeamAOdds       *float32 `json:"team_a_odds"`
	TeamBID         *int64   `json:"team_b_id"`
	TeamBNames      *string  `json:"team_b_names"`
	TeamBOdds       *float32 `json:"team_b_odds"`
	Date            *string  `json:"date"`
	MatchLink       *string  `json:"match_link"`
	Status          *string  `json:"status"`
	Winner          *int64   `json:"winner"`
}
