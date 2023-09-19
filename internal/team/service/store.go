package service

import (
	"context"

	"github.com/x-sports/internal/team"
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

	// CreateTeam creates a new team and return the
	// created team ID.
	CreateTeam(ctx context.Context, team team.Team) (int64, error)

	// GetAllTeams returns all teams.
	GetAllTeams(ctx context.Context) ([]team.Team, error)
}
