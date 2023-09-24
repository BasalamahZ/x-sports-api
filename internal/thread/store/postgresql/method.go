package postgresql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/x-sports/internal/thread"
)

func (sc *storeClient) CreateThread(ctx context.Context, reqThread thread.Thread) (int64, error) {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"title":        reqThread.Title,
		"game_id":      reqThread.GameID,
		"description":  reqThread.Description,
		"image_thread": reqThread.ImageThread,
		"date":         reqThread.Date,
		"create_time":  reqThread.CreateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryCreateThread, argsKV)
	if err != nil {
		return 0, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return 0, err
	}
	query = sc.q.Rebind(query)

	// execute query
	var threadID int64
	err = sc.q.QueryRowx(query, args...).Scan(&threadID)
	if err != nil {
		return 0, err
	}

	return threadID, nil
}

func (sc *storeClient) GetAllThreads(ctx context.Context) ([]thread.Thread, error) {
	query := fmt.Sprintf(queryGetThreads, "")

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

	// read thread
	thread := make([]thread.Thread, 0)
	for rows.Next() {
		var row threadDB
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		thread = append(thread, row.format())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return thread, nil
}

func (sc *storeClient) UpdateThread(ctx context.Context, reqThread thread.Thread) error {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"id":           reqThread.ID,
		"title":        reqThread.Title,
		"game_id":      reqThread.GameID,
		"description":  reqThread.Description,
		"image_thread": reqThread.ImageThread,
		"date":         reqThread.Date,
		"create_time":  reqThread.CreateTime,
		"update_time":  reqThread.UpdateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryUpdateThread, argsKV)
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

func (sc *storeClient) GetThreadByID(ctx context.Context, threadID int64) (thread.Thread, error) {
	query := fmt.Sprintf(queryGetThreads, "WHERE t.id = $1")

	// query single row
	var tdb threadDB
	err := sc.q.QueryRowx(query, threadID).StructScan(&tdb)
	if err != nil {
		return thread.Thread{}, err
	}

	return tdb.format(), nil
}
