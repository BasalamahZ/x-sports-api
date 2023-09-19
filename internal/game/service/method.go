package service

import (
	"context"

	"github.com/x-sports/internal/game"
)

func (s *service) CreateGame(ctx context.Context, reqGame game.Game) (int64, error) {
	// validate field
	err := validateGame(reqGame)
	if err != nil {
		return 0, err
	}

	reqGame.CreateTime = s.timeNow()

	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return 0, err
	}

	gameID, err := pgStoreClient.CreateGame(ctx, reqGame)
	if err != nil {
		return 0, err
	}

	return gameID, nil
}

func (s *service) GetAllGames(ctx context.Context) ([]game.Game, error) {
	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return nil, err
	}

	// get all games from postgre
	games, err := pgStoreClient.GetAllGames(ctx)
	if err != nil {
		return nil, err
	}

	return games, nil
}

// validateGame validates fields of the given Game
// whether its comply the predetermined rules.
func validateGame(reqGame game.Game) error {
	if reqGame.GameNames == "" {
		return game.ErrInvalidGameNames
	}

	if reqGame.GameIcons == "" {
		return game.ErrInvalidGameIcons
	}

	return nil
}
