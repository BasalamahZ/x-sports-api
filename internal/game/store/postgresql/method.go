package postgresql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/x-sports/internal/game"
)

func (sc *storeClient) CreateGame(ctx context.Context, reqGame game.Game) (int64, error) {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"game_names":  reqGame.GameNames,
		"game_icons":  reqGame.GameIcons,
		"create_time": reqGame.CreateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryCreateGame, argsKV)
	if err != nil {
		return 0, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return 0, err
	}
	query = sc.q.Rebind(query)

	// execute query
	var gameID int64
	err = sc.q.QueryRowx(query, args...).Scan(&gameID)
	if err != nil {
		return 0, err
	}

	return gameID, nil
}

func (sc *storeClient) GetAllGames(ctx context.Context) ([]game.Game, error) {
	query := fmt.Sprintf(queryGetGames, "")

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

	// read games
	games := make([]game.Game, 0)
	for rows.Next() {
		var row gameDB
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		games = append(games, row.format())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}

func (sc *storeClient) UpdateGame(ctx context.Context, reqGame game.Game) error {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"id":          reqGame.ID,
		"game_names":  reqGame.GameNames,
		"game_icons":  reqGame.GameIcons,
		"create_time": reqGame.CreateTime,
		"update_time": reqGame.UpdateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryUpdateGame, argsKV)
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

func (sc *storeClient) GetGameByID(ctx context.Context, gameID int64) (game.Game, error) {
	query := fmt.Sprintf(queryGetGames, "WHERE g.id = $1")

	// query single row
	var gdb gameDB
	err := sc.q.QueryRowx(query, gameID).StructScan(&gdb)
	if err != nil {
		return game.Game{}, err
	}

	return gdb.format(), nil
}
