package http

import "github.com/x-sports/internal/game"

// formatGame formats the given game into the
// respective HTTP-format object.
func formatGame(g game.Game) (gameHTTP, error) {
	return gameHTTP{
		ID:        &g.ID,
		GameNames: &g.GameNames,
		GameIcons: &g.GameIcons,
	}, nil
}
