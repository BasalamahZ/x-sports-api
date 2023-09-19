package team

import "errors"

var (
	// ErrInvalidTeamNames is returned when the given team names is
	// invalid.
	ErrInvalidTeamNames = errors.New("invalid team names")

	// ErrInvalidGameID is returned when the given game id is
	// invalid.
	ErrInvalidGameID = errors.New("invalid game id")
)
