package match

import (
	"context"
	"time"
)

type Service interface {
	// CreateMatch creates a new match and return the
	// created Match ID.
	CreateMatch(ctx context.Context, match Match) (int64, error)

	// GetAllMatchs returns all matchs and filter by game id.
	GetAllMatchs(ctx context.Context, gameID int64, status Status) ([]Match, error)

	// GetMatchByID returns a match with the given
	// match ID.
	GetMatchByID(ctx context.Context, matchID int64) (Match, error)

	// UpdateMatch updates existing match with the given
	// match data.
	//
	// UpdateMatch do updates on all main attributes
	// except ID, CreateBy, and CreateTime. So, make sure to
	// use current values in the given data if do not want
	// to update some specific attributes.
	UpdateMatch(ctx context.Context, match Match) error
}

type Match struct {
	ID              int64
	BlockChainID    uint64
	TournamentNames string
	GameID          int64
	GameNames       string // derived
	GameIcons       string // derived
	TeamAID         int64
	TeamANames      string // derived
	TeamAIcons      string // derived
	TeamAOdds       float32
	TeamBID         int64
	TeamBNames      string // derived
	TeamBIcons      string // derived
	TeamBOdds       float32
	Date            time.Time
	MatchLink       string
	Status          Status
	Winner          int64
	CreateTime      time.Time
	UpdateTime      time.Time
}

// Type denotes type of a status.
type Status int

// Followings are the known status.
const (
	StatusUnknown   Status = 0
	StatusUpcoming  Status = 1
	StatusOngoing   Status = 2
	StatusCompleted Status = 3
)

var (
	// StatusList is a list of valid status levels.
	StatusList = map[Status]struct{}{
		StatusUpcoming:  {},
		StatusOngoing:   {},
		StatusCompleted: {},
	}

	// StatusName maps status level to it's string
	// representation.
	statusName = map[Status]string{
		StatusUpcoming:  "upcoming",
		StatusOngoing:   "ongoing",
		StatusCompleted: "completed",
	}
)

// Value returns int value of a status type.
func (s Status) Value() int {
	return int(s)
}

// String returns string representaion of a status type.
func (s Status) String() string {
	return statusName[s]
}
