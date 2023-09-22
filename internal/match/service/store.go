package service

import (
	"context"

	"github.com/x-sports/internal/match"
)

// PGStore is the PostgreSQL store for admin service.
type PGStore interface {
	NewClient(useTx bool) (PGStoreClient, error)
}

type PGStoreClient interface {
	// Commit commits the transaction.
	Commit() error

	// Rollback aborts the transaction.
	Rollback() error

	// CreateMatch creates a new match and return the
	// created Match ID.
	CreateMatch(ctx context.Context, match match.Match) (int64, error)

	// GetAllMatchs returns all matchs and filter by game id.
	GetAllMatchs(ctx context.Context, gameID int64, status match.Status) ([]match.Match, error)

	// GetMatchByID returns a match with the given
	// match ID.
	GetMatchByID(ctx context.Context, matchID int64) (match.Match, error)

	// UpdateMatch updates existing match with the given
	// match data.
	//
	// UpdateMatch do updates on all main attributes
	// except ID, CreateBy, and CreateTime. So, make sure to
	// use current values in the given data if do not want
	// to update some specific attributes.
	UpdateMatch(ctx context.Context, match match.Match) error
}
