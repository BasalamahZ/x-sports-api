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

	// these value should be same for all users
	var (
		createTime = s.timeNow()
	)

	reqGame.CreateTime = createTime

	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(true)
	if err != nil {
		return 0, err
	}

	gameID, err := pgStoreClient.CreateGame(ctx, reqGame)
	if err != nil {
		return 0, err
	}

	return gameID, nil
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
