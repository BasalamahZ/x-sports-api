package http

import (
	"errors"

	"github.com/x-sports/internal/thread"
)

// Followings are the known errors from Thread HTTP handlers.
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

	// errInvalidTitle is returned when the given title is
	// invalid.
	errInvalidTitle = errors.New("INVALID_TITLE")

	// errInvalidThreadID is returned when the given thread ID is
	// invalid.
	errInvalidThreadID = errors.New("INVALID_THREAD_ID")

	// errInvalidGameID is returned when the given game ID is
	// invalid.
	errInvalidGameID = errors.New("INVALID_GAME_ID")

	// errInvalidDescription is returned when the given
	// description is invalid.
	errInvalidDescription = errors.New("INVALID_DESCRIPTION")

	// errInvalidDate is returned when the given thread
	// date is invalid.
	errInvalidDate = errors.New("INVALID_DATE")

	// errInvalidImageThread is returned when the given
	// image thread is invalid.
	errInvalidImageThread = errors.New("INVALID_IMAGE_THREAD")

	// errInvalidTimeFormat is returned when the given time
	// string format is invalid.
	errInvalidTimeFormat = errors.New("INVALID_TIME_FORMAT")

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
		thread.ErrInvalidTitle:       errInvalidTitle,
		thread.ErrInvalidThreadID:      errInvalidThreadID,
		thread.ErrInvalidGameID:      errInvalidGameID,
		thread.ErrInvalidDescription: errInvalidDescription,
		thread.ErrInvalidDate:        errInvalidDate,
		thread.ErrInvalidImageThread:   errInvalidImageThread,
	}
)
