package http

import (
	"errors"

	"github.com/x-sports/internal/team"
)

// Followings are the known errors from Game HTTP handlers.
var (
	// errBadRequest is returned when the given request is
	// bad/invalid.
	errBadRequest = errors.New("BAD_REQUEST")

	// errInternalServer is returned when there is an
	// unexpected error encountered when processing a request.
	errInternalServer = errors.New("INTERNAL_SERVER_ERROR")

	// errInvalidToken is returned when the given token is
	// invalid.
	errInvalidToken = errors.New("INVALID_TOKEN")

	// errInvalidTeamNames is returned when the given team names is
	// invalid.
	errInvalidTeamNames = errors.New("INVALID_TEAM_NAMES")

	// errInvalidGameID is returned when the given game ID is
	// invalid.
	errInvalidGameID = errors.New("INVALID_GAME_ID")

	// errMethodNotAllowed is returned when accessing not
	// allowed HTTP method.
	errMethodNotAllowed = errors.New("METHOD_NOT_ALLOWED")

	// errRequestTimeout is returned when processing time
	// has reached the timeout limit.
	errRequestTimeout = errors.New("REQUEST_TIMEOUT")

	// errUnauthorizedAccess is returned when the request
	// is unaothorized.
	errUnauthorizedAccess = errors.New("UNAUTHORIZED_ACCESS")
)

var (
	// mapHTTPError maps service error into HTTP error that
	// categorize as bad request error.
	//
	// Internal server error-related should not be mapped
	// here, and the handler should just return `errInternal`
	// as the error instead
	mapHTTPError = map[error]error{
		team.ErrInvalidTeamNames: errInvalidTeamNames,
		team.ErrInvalidGameID:    errInvalidGameID,
	}
)
