package team

import (
	"context"
	"time"
)

type Service interface {
	// CreateTeam creates a new team and return the
	// created team ID.
	CreateTeam(ctx context.Context, team Team) (int64, error)

	// GetAllTeams returns all teams.
	GetAllTeams(ctx context.Context) ([]Team, error)
}

type Team struct {
	ID         int64
	TeamNames  string
	TeamIcons  string
	GameID     int64
	GameNames  string // derived
	GameIcons  string // derived
	CreateTime time.Time
	UpdateTime time.Time
}
