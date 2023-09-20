package http

import (
	"github.com/x-sports/internal/team"
)

// formatTeam formats the given team into the
// respective HTTP-format object.
func formatTeam(t team.Team) (teamHTTP, error) {
	return teamHTTP{
		ID:        &t.ID,
		TeamNames: &t.TeamNames,
		GameID:    &t.GameID,
		GameNames: &t.GameNames,
		GameIcons: &t.GameIcons,
	}, nil
}
