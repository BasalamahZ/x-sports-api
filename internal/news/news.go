package news

import (
	"context"
	"time"
)

type Service interface {
	// CreateNews creates a new news and return the
	// created news ID.
	CreateNews(ctx context.Context, news News) (int64, error)

	// GetAllNews returns all news and filter by game id.
	GetAllNews(ctx context.Context, gameID int64) ([]News, error)

	// GetNewsByID returns a news with the given
	// news ID.
	GetNewsByID(ctx context.Context, newsID int64) (News, error)

	// UpdateNews updates existing news with the given
	// news data.
	//
	// UpdateNews do updates on all main attributes
	// except ID, CreateBy, and CreateTime. So, make sure to
	// use current values in the given data if do not want
	// to update some specific attributes.
	UpdateNews(ctx context.Context, news News) error
}

type News struct {
	ID          int64
	Title       string
	GameID      int64
	GameNames   string // derived
	GameIcons   string // derived
	Description string
	ImageNews   string
	Date        time.Time
	CreateTime  time.Time
	UpdateTime  time.Time
}
