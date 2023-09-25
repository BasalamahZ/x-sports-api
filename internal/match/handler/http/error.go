package http

import (
	"errors"

	"github.com/x-sports/internal/match"
)

// Followings are the known errors from Match HTTP handlers.
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

	// errInvalidTournamentNames is returned when the given tournament names is
	// invalid.
	errInvalidTournamentNames = errors.New("INVALID_TOURNAMENT_NAMES")

	// errInvalidMatchID is returned when the given match ID is
	// invalid.
	errInvalidMatchID = errors.New("INVALID_MATCH_ID")

	// errInvalidBlockChainID is returned when the given blockchain ID is
	// invalid.
	errInvalidBlockChainID = errors.New("INVALID_BLOCKCHAIN_ID")

	// errInvalidGameID is returned when the given game ID is
	// invalid.
	errInvalidGameID = errors.New("INVALID_GAME_ID")

	// errInvalidStatus is returned when the given match
	// status is invalid.
	errInvalidStatus = errors.New("INVALID_STATUS")

	// errInvalidDate is returned when the given match
	// date is invalid.
	errInvalidDate = errors.New("INVALID_DATE")

	// errInvalidTeamID is returned when the given
	// team ID is invalid.
	errInvalidTeamID = errors.New("INVALID_TEAM_ID")

	// errInvalidTeamOdds is returned when the given
	// team odds is invalid.
	errInvalidTeamOdds = errors.New("INVALID_TEAM_ODDS")

	// errInvalidWinner is returned when the given
	// winner is invalid.
	errInvalidWinner = errors.New("INVALID_WINNER")

	// errInvalidMatchLink is returned when the given
	// match link is invalid.
	errInvalidMatchLink = errors.New("INVALID_MATCH_LINK")

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
		match.ErrInvalidTournamentNames: errInvalidTournamentNames,
		match.ErrInvalidBlockChainID:    errInvalidBlockChainID,
		match.ErrInvalidMatchID:         errInvalidMatchID,
		match.ErrInvalidGameID:          errInvalidGameID,
		match.ErrInvalidTeamID:          errInvalidTeamID,
		match.ErrInvalidTeamOdds:        errInvalidTeamOdds,
		match.ErrInvalidStatus:          errInvalidStatus,
		match.ErrInvalidDate:            errInvalidDate,
		match.ErrInvalidWinner:          errInvalidWinner,
		match.ErrInvalidMatchLink:       errInvalidMatchLink,
	}
)
