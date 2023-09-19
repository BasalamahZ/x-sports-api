package service

import (
	"context"

	"github.com/x-sports/internal/game"
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

	// CreateGame creates a new game and return the
	// created game ID.
	CreateGame(ctx context.Context, game game.Game) (int64, error)

	// GetAllGames returns all games.
	GetAllGames(ctx context.Context) ([]game.Game, error)
}
