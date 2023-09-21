package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/x-sports/internal/news"
)

func (sc *storeClient) CreateNews(ctx context.Context, reqNews news.News) (int64, error) {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"title":       reqNews.Title,
		"game_id":     reqNews.GameID,
		"description": reqNews.Description,
		"image_news":  reqNews.ImageNews,
		"date":        reqNews.Date,
		"create_time": reqNews.CreateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryCreateNews, argsKV)
	if err != nil {
		return 0, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return 0, err
	}
	query = sc.q.Rebind(query)

	// execute query
	var newsID int64
	err = sc.q.QueryRowx(query, args...).Scan(&newsID)
	if err != nil {
		return 0, err
	}

	return newsID, nil
}

func (sc *storeClient) GetAllNews(ctx context.Context, gameID int64) ([]news.News, error) {
	// define variables to custom query
	argsKV := make(map[string]interface{})
	addConditions := make([]string, 0)

	if gameID > 0 {
		addConditions = append(addConditions, "n.game_id = :game_id")
		argsKV["game_id"] = gameID
	}

	// construct strings to custom query
	addCondition := strings.Join(addConditions, " AND ")

	// since the query does not contains "WHERE" yet, need
	// to add it if needed
	if len(addConditions) > 0 {
		addCondition = fmt.Sprintf("WHERE %s", addCondition)
	}
	query := fmt.Sprintf(queryGetNews, addCondition)

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

	// read news
	news := make([]news.News, 0)
	for rows.Next() {
		var row newsDB
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		news = append(news, row.format())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return news, nil
}

func (sc *storeClient) UpdateNews(ctx context.Context, reqNews news.News) error {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"id":          reqNews.ID,
		"title":       reqNews.Title,
		"game_id":     reqNews.GameID,
		"description": reqNews.Description,
		"image_news":  reqNews.ImageNews,
		"date":        reqNews.Date,
		"create_time": reqNews.CreateTime,
		"update_time": reqNews.UpdateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryUpdateNews, argsKV)
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

func (sc *storeClient) GetNewsByID(ctx context.Context, newsID int64) (news.News, error) {
	query := fmt.Sprintf(queryGetNews, "WHERE n.id = $1")

	// query single row
	var ndb newsDB
	err := sc.q.QueryRowx(query, newsID).StructScan(&ndb)
	if err != nil {
		return news.News{}, err
	}

	return ndb.format(), nil
}
