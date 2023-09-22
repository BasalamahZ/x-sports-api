package postgresql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/x-sports/internal/team"
)

func (sc *storeClient) CreateTeam(ctx context.Context, reqTeam team.Team) (int64, error) {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"team_names":  reqTeam.TeamNames,
		"team_icons":  reqTeam.TeamIcons,
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
	query := fmt.Sprintf(queryGetTeams, "")

	// prepare query
	query, args, err := sqlx.Named(query, map[string]interface{}{})
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

func (sc *storeClient) UpdateTeam(ctx context.Context, reqTeam team.Team) error {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"id":          reqTeam.ID,
		"team_names":  reqTeam.TeamNames,
		"team_icons":  reqTeam.TeamIcons,
		"game_id":     reqTeam.GameID,
		"create_time": reqTeam.CreateTime,
		"update_time": reqTeam.UpdateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryUpdateTeam, argsKV)
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

func (sc *storeClient) GetTeamByID(ctx context.Context, teamID int64) (team.Team, error) {
	query := fmt.Sprintf(queryGetTeams, "WHERE t.id = $1")

	// query single row
	var tdb teamDB
	err := sc.q.QueryRowx(query, teamID).StructScan(&tdb)
	if err != nil {
		return team.Team{}, err
	}

	return tdb.format(), nil
}
