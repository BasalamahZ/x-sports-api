package http

import (
	"github.com/x-sports/internal/match"
)

// formatMatch formats the given match into the
// respective HTTP-format object.
func formatMatch(m match.Match) (matchHTTP, error) {
	date := m.Date.Format(dateFormat)
	status := m.Status.String()

	return matchHTTP{
		ID:              &m.ID,
		BlockChainID:    &m.BlockChainID,
		TournamentNames: &m.TournamentNames,
		GameID:          &m.GameID,
		GameNames:       &m.GameNames,
		GameIcons:       &m.GameIcons,
		TeamAID:         &m.TeamAID,
		TeamANames:      &m.TeamANames,
		TeamAIcons:      &m.TeamAIcons,
		TeamAOdds:       &m.TeamAOdds,
		TeamBID:         &m.TeamBID,
		TeamBNames:      &m.TeamBNames,
		TeamBIcons:      &m.TeamBIcons,
		TeamBOdds:       &m.TeamBOdds,
		Date:            &date,
		Status:          &status,
		MatchLink:       &m.MatchLink,
		Winner:          &m.Winner,
	}, nil
}

// parseStatus returns match.Status from the given string.
func parseStatus(req string) (match.Status, error) {
	switch req {
	case match.StatusUpcoming.String():
		return match.StatusUpcoming, nil
	case match.StatusOngoing.String():
		return match.StatusOngoing, nil
	case match.StatusCompleted.String():
		return match.StatusCompleted, nil
	}
	return match.StatusUnknown, errInvalidStatus
}
