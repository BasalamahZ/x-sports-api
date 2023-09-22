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

	// GetGameByID returns a game with the given
	// game ID.
	GetGameByID(ctx context.Context, gameID int64) (game.Game, error)

	// UpdateGame updates existing game with the given
	// game data.
	//
	// UpdateGame do updates on all main attributes
	// except ID, CreateBy, and CreateTime. So, make sure to
	// use current values in the given data if do not want
	// to update some specific attributes.
	UpdateGame(ctx context.Context, game game.Game) error
}
