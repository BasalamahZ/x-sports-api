package service

import (
	"context"

	"github.com/x-sports/internal/match"
)

func (s *service) CreateMatch(ctx context.Context, reqMatch match.Match) (int64, error) {
	// validate field
	err := validateMatch(reqMatch)
	if err != nil {
		return 0, err
	}

	reqMatch.CreateTime = s.timeNow()

	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return 0, err
	}

	matchID, err := pgStoreClient.CreateMatch(ctx, reqMatch)
	if err != nil {
		return 0, err
	}

	return matchID, nil
}

func (s *service) GetAllMatchs(ctx context.Context, gameID int64, status match.Status) ([]match.Match, error) {
	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return nil, err
	}

	// get all match from postgre
	match, err := pgStoreClient.GetAllMatchs(ctx, gameID, status)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func (s *service) UpdateMatch(ctx context.Context, reqMatch match.Match) error {
	// validate field
	err := validateMatch(reqMatch)
	if err != nil {
		return err
	}

	// modify fields
	reqMatch.UpdateTime = s.timeNow()

	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	if reqMatch.Status == match.StatusCompleted {
		if (reqMatch.Winner != reqMatch.TeamAID && reqMatch.Winner != reqMatch.TeamBID) || reqMatch.Winner <= 0 {
			return match.ErrInvalidWinner
		}
	}

	// updates match in pgstore
	err = pgStoreClient.UpdateMatch(ctx, reqMatch)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetMatchByID(ctx context.Context, matchID int64) (match.Match, error) {
	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return match.Match{}, err
	}

	// get match from pgstore
	result, err := pgStoreClient.GetMatchByID(ctx, matchID)
	if err != nil {
		return match.Match{}, err
	}

	return result, nil
}

// validateMatch validates fields of the given Match
// whether its comply the predetermined rules.
func validateMatch(reqMatch match.Match) error {
	if reqMatch.TournamentNames == "" {
		return match.ErrInvalidTournamentNames
	}

	if reqMatch.GameID <= 0 {
		return match.ErrInvalidGameID
	}

	if reqMatch.TeamAID <= 0 {
		return match.ErrInvalidTeamID
	}

	if reqMatch.TeamBID <= 0 {
		return match.ErrInvalidTeamID
	}

	if reqMatch.TeamAOdds <= 0 {
		return match.ErrInvalidTeamOdds
	}

	if reqMatch.TeamBOdds <= 0 {
		return match.ErrInvalidTeamOdds
	}

	if reqMatch.Date.IsZero() {
		return match.ErrInvalidDate
	}

	if _, valid := match.StatusList[reqMatch.Status]; !valid {
		return match.ErrInvalidStatus
	}

	if reqMatch.MatchLink == "" {
		return match.ErrInvalidMatchLink
	}

	return nil
}
