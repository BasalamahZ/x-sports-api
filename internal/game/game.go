package game

import (
	"context"
	"time"
)

type Service interface {
	// CreateGame creates a new game and return the
	// created game ID.
	CreateGame(ctx context.Context, game Game) (int64, error)

	// GetAllGames returns all games.
	GetAllGames(ctx context.Context) ([]Game, error)
}

type Game struct {
	ID         int64
	GameNames  string
	GameIcons  string
	CreateTime time.Time
	UpdateTime time.Time
}
