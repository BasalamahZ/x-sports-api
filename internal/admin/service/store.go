package service

import (
	"context"

	"github.com/x-sports/internal/admin"
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

	// GetUserByEmail selects a user with the given
	// email.
	GetUserByEmail(ctx context.Context, email string) (admin.Admin, error)
}
