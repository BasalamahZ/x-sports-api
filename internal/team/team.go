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
	
	// GetTeamsByID returns a team with the given
	// team ID.
	GetTeamByID(ctx context.Context, teamID int64) (Team, error)

	// UpdateTeam updates existing team with the given
	// team data.
	//
	// Updateteams do updates on all main attributes
	// except ID, CreateBy, and CreateTime. So, make sure to
	// use current values in the given data if do not want
	// to update some specific attributes.
	UpdateTeam(ctx context.Context, team Team) error
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
