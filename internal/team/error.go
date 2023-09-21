package team

import "errors"

var (
	// ErrInvalidTeamNames is returned when the given team names is
	// invalid.
	ErrInvalidTeamNames = errors.New("invalid team names")

	// ErrInvalidTeamIcons is returned when the given team icons is
	// invalid.
	ErrInvalidTeamIcons = errors.New("invalid team icons")

	// ErrInvalidGameID is returned when the given game id is
	// invalid.
	ErrInvalidGameID = errors.New("invalid game id")
)
