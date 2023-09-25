package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/x-sports/internal/match"
)

func (sc *storeClient) CreateMatch(ctx context.Context, reqMatch match.Match) (int64, error) {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"tournament_names": reqMatch.TournamentNames,
		"blockchain_id":    reqMatch.BlockChainID,
		"game_id":          reqMatch.GameID,
		"game_names":       reqMatch.GameNames,
		"game_icons":       reqMatch.GameIcons,
		"team_a_id":        reqMatch.TeamAID,
		"team_a_names":     reqMatch.TeamANames,
		"team_a_icons":     reqMatch.TeamAIcons,
		"team_a_odds":      reqMatch.TeamAOdds,
		"team_b_id":        reqMatch.TeamBID,
		"team_b_names":     reqMatch.TeamBNames,
		"team_b_icons":     reqMatch.TeamBIcons,
		"team_b_odds":      reqMatch.TeamBOdds,
		"date":             reqMatch.Date,
		"match_link":       reqMatch.MatchLink,
		"status":           reqMatch.Status,
		"create_time":      reqMatch.CreateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryCreateMatch, argsKV)
	if err != nil {
		return 0, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return 0, err
	}
	query = sc.q.Rebind(query)

	// execute query
	var matchID int64
	err = sc.q.QueryRowx(query, args...).Scan(&matchID)
	if err != nil {
		return 0, err
	}

	return matchID, nil
}

func (sc *storeClient) GetAllMatchs(ctx context.Context, gameID int64, status match.Status) ([]match.Match, error) {
	// define variables to custom query
	argsKV := make(map[string]interface{})
	addConditions := make([]string, 0)

	if gameID > 0 {
		addConditions = append(addConditions, "m.game_id = :game_id")
		argsKV["game_id"] = gameID
	}

	if status > 0 {
		addConditions = append(addConditions, "m.status = :status")
		argsKV["status"] = status
	}

	// construct strings to custom query
	addCondition := strings.Join(addConditions, " AND ")

	// since the query does not contains "WHERE" yet, need
	// to add it if needed
	if len(addConditions) > 0 {
		addCondition = fmt.Sprintf("WHERE %s", addCondition)
	}
	query := fmt.Sprintf(queryGetMatchs, addCondition)

	// prepare query
	query, args, err := sqlx.Named(query, argsKV)
	if err != nil {
		return nil, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = sc.q.Rebind(query)

	// query to database
	rows, err := sc.q.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// read match
	match := make([]match.Match, 0)
	for rows.Next() {
		var row matchDB
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		match = append(match, row.format())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return match, nil
}

func (sc *storeClient) UpdateMatch(ctx context.Context, reqMatch match.Match) error {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"id":               reqMatch.ID,
		"blockchain_id":    reqMatch.BlockChainID,
		"tournament_names": reqMatch.TournamentNames,
		"game_id":          reqMatch.GameID,
		"game_names":       reqMatch.GameNames,
		"game_icons":       reqMatch.GameIcons,
		"team_a_id":        reqMatch.TeamAID,
		"team_a_names":     reqMatch.TeamANames,
		"team_a_icons":     reqMatch.TeamAIcons,
		"team_a_odds":      reqMatch.TeamAOdds,
		"team_b_id":        reqMatch.TeamBID,
		"team_b_names":     reqMatch.TeamBNames,
		"team_b_icons":     reqMatch.TeamBIcons,
		"team_b_odds":      reqMatch.TeamBOdds,
		"date":             reqMatch.Date,
		"match_link":       reqMatch.MatchLink,
		"status":           reqMatch.Status,
		"winner":           reqMatch.Winner,
		"create_time":      reqMatch.CreateTime,
		"update_time":      reqMatch.UpdateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryUpdateMatch, argsKV)
	if err != nil {
		return err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}
	query = sc.q.Rebind(query)

	// execute query
	_, err = sc.q.Exec(query, args...)
	return err
}

func (sc *storeClient) GetMatchByID(ctx context.Context, matchID int64) (match.Match, error) {
	query := fmt.Sprintf(queryGetMatchs, "WHERE m.id = $1")

	// query single row
	var mdb matchDB
	err := sc.q.QueryRowx(query, matchID).StructScan(&mdb)
	if err != nil {
		return match.Match{}, err
	}

	return mdb.format(), nil
}
