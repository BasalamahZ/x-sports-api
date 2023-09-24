package service

import (
	"context"

	"github.com/x-sports/internal/thread"
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

	// CreateThread creates a new thread and return the
	// created thread ID.
	CreateThread(ctx context.Context, thread thread.Thread) (int64, error)

	// GetAllThreads returns all thread.
	GetAllThreads(ctx context.Context) ([]thread.Thread, error)

	// GetThreadByID returns a thread with the given
	// thread ID.
	GetThreadByID(ctx context.Context, threadID int64) (thread.Thread, error)

	// UpdateThread updates existing thread with the given
	// thread data.
	//
	// Updatethread do updates on all main attributes
	// except ID, CreateBy, and CreateTime. So, make sure to
	// use current values in the given data if do not want
	// to update some specific attributes.
	UpdateThread(ctx context.Context, thread thread.Thread) error
}
