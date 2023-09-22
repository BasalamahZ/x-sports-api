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

func (s *service) UpdateTeam(ctx context.Context, reqTeam team.Team) error {
	// validate field
	err := validateTeam(reqTeam)
	if err != nil {
		return err
	}

	// modify fields
	reqTeam.UpdateTime = s.timeNow()

	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// updates team in pgstore
	err = pgStoreClient.UpdateTeam(ctx, reqTeam)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetTeamByID(ctx context.Context, teamID int64) (team.Team, error) {
	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return team.Team{}, err
	}

	// get Team from pgstore
	result, err := pgStoreClient.GetTeamByID(ctx, teamID)
	if err != nil {
		return team.Team{}, err
	}

	return result, nil
}

// validateTeam validates fields of the given Team
// whether its comply the predetermined rules.
func validateTeam(reqTeam team.Team) error {
	if reqTeam.TeamNames == "" {
		return team.ErrInvalidTeamNames
	}

	if reqTeam.TeamIcons == "" {
		return team.ErrInvalidTeamIcons
	}

	if reqTeam.GameID <= 0 {
		return team.ErrInvalidGameID
	}

	return nil
}
