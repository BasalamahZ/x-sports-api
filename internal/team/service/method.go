package service

import (
	"context"

	"github.com/x-sports/internal/team"
)

func (s *service) CreateTeam(ctx context.Context, reqTeam team.Team) (int64, error) {
	// validate field
	err := validateTeam(reqTeam)
	if err != nil {
		return 0, err
	}

	reqTeam.CreateTime = s.timeNow()

	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return 0, err
	}

	teamID, err := pgStoreClient.CreateTeam(ctx, reqTeam)
	if err != nil {
		return 0, err
	}

	return teamID, nil
}

func (s *service) GetAllTeams(ctx context.Context) ([]team.Team, error) {
	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return nil, err
	}

	// get all teams from postgre
	teams, err := pgStoreClient.GetAllTeams(ctx)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

// validateTeam validates fields of the given Team
// whether its comply the predetermined rules.
func validateTeam(reqTeam team.Team) error {
	if reqTeam.TeamNames == "" {
		return team.ErrInvalidTeamNames
	}

	if reqTeam.GameID <= 0 {
		return team.ErrInvalidGameID
	}

	return nil
}
