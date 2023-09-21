package news

import "errors"

var (
	// ErrInvalidTitle is returned when the given title is
	// invalid.
	ErrInvalidTitle = errors.New("invalid title")

	// ErrInvalidNewsID is returned when the given news id is
	// invalid.
	ErrInvalidNewsID = errors.New("invalid news id")

	// ErrInvalidGameID is returned when the given game id is
	// invalid.
	ErrInvalidGameID = errors.New("invalid game id")

	// ErrInvalidDescription is returned when the given description is
	// invalid.
	ErrInvalidDescription = errors.New("invalid description")

	// ErrInvalidImageNews is returned when the given image news
	// is invalid.
	ErrInvalidImageNews = errors.New("invalid image news")

	// ErrInvalidDate is returned whrn the given date
	// is invalid.
	ErrInvalidDate = errors.New("invalid date")
)
