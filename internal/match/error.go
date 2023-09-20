package match

import "errors"

var (
	// ErrInvalidTournamentNames is returned when the given tournament names is
	// invalid.
	ErrInvalidTournamentNames = errors.New("invalid tournament names")

	// ErrInvalidMatchID is returned when the given match id is
	// invalid.
	ErrInvalidMatchID = errors.New("invalid match id")

	// ErrInvalidGameID is returned when the given game id is
	// invalid.
	ErrInvalidGameID = errors.New("invalid game id")

	// ErrInvalidTeamID is returned when the given team id is
	// invalid.
	ErrInvalidTeamID = errors.New("invalid team id")

	// ErrInvalidTeamOdds is returned when the given team odds is
	// invalid.
	ErrInvalidTeamOdds = errors.New("invalid team odds")

	// ErrInvalidStatus is returned when the given status
	// is invalid.
	ErrInvalidStatus = errors.New("invalid status")

	// ErrInvalidDate is returned whrn the given date
	// is invalid.
	ErrInvalidDate = errors.New("invalid date")

	// ErrInvalidWinner is returned whrn the given winner
	// is invalid.
	ErrInvalidWinner = errors.New("invalid winner")

	// ErrInvalidMatchLink is returned whrn the given match link
	// is invalid.
	ErrInvalidMatchLink = errors.New("invalid match link")
)
