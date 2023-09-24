package thread

import "errors"

var (
	// ErrInvalidTitle is returned when the given title is
	// invalid.
	ErrInvalidTitle = errors.New("invalid title")

	// ErrInvalidThreadID is returned when the given thread id is
	// invalid.
	ErrInvalidThreadID = errors.New("invalid thread id")

	// ErrInvalidGameID is returned when the given game id is
	// invalid.
	ErrInvalidGameID = errors.New("invalid game id")

	// ErrInvalidDescription is returned when the given description is
	// invalid.
	ErrInvalidDescription = errors.New("invalid description")

	// ErrInvalidImageThread is returned when the given image thread
	// is invalid.
	ErrInvalidImageThread = errors.New("invalid image thread")

	// ErrInvalidDate is returned whrn the given date
	// is invalid.
	ErrInvalidDate = errors.New("invalid date")
)
