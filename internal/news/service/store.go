package service

import (
	"context"

	"github.com/x-sports/internal/news"
)

// PGStore is the PostgreSQL store for admin service.
type PGStore interface {
	NewClient(useTx bool) (PGStoreClient, error)
}

type PGStoreClient interface {
	// Commit commits the transaction.
	Commit() error

	// Rollback aborts the transaction.
	Rollback() error

	// CreateNews creates a new news and return the
	// created news ID.
	CreateNews(ctx context.Context, news news.News) (int64, error)

	// GetAllNews returns all news.
	GetAllNews(ctx context.Context, gameID int64) ([]news.News, error)

	// GetNewsByID returns a news with the given
	// news ID.
	GetNewsByID(ctx context.Context, newsID int64) (news.News, error)

	// UpdateNews updates existing news with the given
	// news data.
	//
	// UpdateNews do updates on all main attributes
	// except ID, CreateBy, and CreateTime. So, make sure to
	// use current values in the given data if do not want
	// to update some specific attributes.
	UpdateNews(ctx context.Context, news news.News) error
}
