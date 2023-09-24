package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/x-sports/internal/admin"
	"github.com/x-sports/internal/thread"
)

var (
	errUnknownConfig = errors.New("unknown config name")
)

// dateFormat denotes the standard date format used in
// thread HTTP request and response.
var dateFormat = "2 January 2006"

// Handler contains admin HTTP-handlers.
type Handler struct {
	handlers map[string]*handler
	thread   thread.Service
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
	// HandlerThread denotes HTTP handler to interact
	// with a thread
	HandlerThread = HandlerIdentity{
		Name: "thread",
		URL:  "/threads/{id}",
	}

	// HandlerThreads denotes HTTP handler to interact
	// with a threads
	HandlerThreads = HandlerIdentity{
		Name: "threads",
		URL:  "/threads",
	}
)

// New creates a new Handler.
func New(thread thread.Service, admin admin.Service, identities []HandlerIdentity) (*Handler, error) {
	h := &Handler{
		handlers: make(map[string]*handler),
		thread:   thread,
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
	case HandlerThread.Name:
		httpHandler = &threadHandler{
			thread: h.thread,
			admin:  h.admin,
		}
	case HandlerThreads.Name:
		httpHandler = &threadsHandler{
			thread: h.thread,
			admin:  h.admin,
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

// threadHTTP denotes user object in HTTP request or response
// body.
type threadHTTP struct {
	ID          *int64  `json:"id"`
	Title       *string `json:"title"`
	GameID      *int64  `json:"game_id"`
	GameNames   *string `json:"game_names"`
	GameIcons   *string `json:"game_icons"`
	Description *string `json:"description"`
	ImageThread *string `json:"image_thread"`
	Date        *string `json:"date"`
}
