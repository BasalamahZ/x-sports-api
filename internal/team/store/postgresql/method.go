package postgresql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/x-sports/internal/team"
)

func (sc *storeClient) CreateTeam(ctx context.Context, reqTeam team.Team) (int64, error) {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"team_names":  reqTeam.TeamNames,
		"game_id":     reqTeam.GameID,
		"create_time": reqTeam.CreateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryCreateTeam, argsKV)
	if err != nil {
		return 0, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return 0, err
	}
	query = sc.q.Rebind(query)

	// execute query
	var teamID int64
	err = sc.q.QueryRowx(query, args...).Scan(&teamID)
	if err != nil {
		return 0, err
	}

	return teamID, nil
}

func (sc *storeClient) GetAllTeams(ctx context.Context) ([]team.Team, error) {
	// prepare query
	query, args, err := sqlx.Named(queryGetTeams, map[string]interface{}{})
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

	// read teams
	teams := make([]team.Team, 0)
	for rows.Next() {
		var row teamDB
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		teams = append(teams, row.format())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}
