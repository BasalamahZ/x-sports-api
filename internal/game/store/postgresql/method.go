package postgresql

import (
	"context"

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
