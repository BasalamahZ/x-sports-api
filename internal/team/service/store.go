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

	// GetTeamsByID returns a team with the given
	// team ID.
	GetTeamByID(ctx context.Context, teamID int64) (team.Team, error)

	// UpdateTeam updates existing team with the given
	// team data.
	//
	// Updateteams do updates on all main attributes
	// except ID, CreateBy, and CreateTime. So, make sure to
	// use current values in the given data if do not want
	// to update some specific attributes.
	UpdateTeam(ctx context.Context, team team.Team) error
}
