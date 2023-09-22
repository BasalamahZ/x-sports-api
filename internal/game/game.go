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

	// GetGameByID returns a game with the given
	// game ID.
	GetGameByID(ctx context.Context, gameID int64) (Game, error)

	// UpdateGame updates existing game with the given
	// game data.
	//
	// UpdateGame do updates on all main attributes
	// except ID, CreateBy, and CreateTime. So, make sure to
	// use current values in the given data if do not want
	// to update some specific attributes.
	UpdateGame(ctx context.Context, game Game) error
}

type Game struct {
	ID         int64
	GameNames  string
	GameIcons  string
	CreateTime time.Time
	UpdateTime time.Time
}
