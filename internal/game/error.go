package game

import "errors"

var (
	// ErrInvalidGamesNames is returned when the given game names is
	// invalid.
	ErrInvalidGameNames = errors.New("invalid game names")

	// ErrInvalidGameIcons is returned when the given game icons is
	// invalid.
	ErrInvalidGameIcons = errors.New("invalid game icons")

	// ErrInvalidGameID is returned when the given game id is
	// invalid.
	ErrInvalidGameID = errors.New("invalid game id")
)
