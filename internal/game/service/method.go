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

func (s *service) UpdateGame(ctx context.Context, reqGame game.Game) error {
	// validate field
	err := validateGame(reqGame)
	if err != nil {
		return err
	}

	// modify fields
	reqGame.UpdateTime = s.timeNow()

	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// updates Game in pgstore
	err = pgStoreClient.UpdateGame(ctx, reqGame)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetGameByID(ctx context.Context, gameID int64) (game.Game, error) {
	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return game.Game{}, err
	}

	// get game from pgstore
	result, err := pgStoreClient.GetGameByID(ctx, gameID)
	if err != nil {
		return game.Game{}, err
	}

	return result, nil
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
