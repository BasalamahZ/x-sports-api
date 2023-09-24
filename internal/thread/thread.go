package thread

import (
	"context"
	"time"
)

type Service interface {
	// CreateThread creates a new thread and return the
	// created thread ID.
	CreateThread(ctx context.Context, thread Thread) (int64, error)

	// GetAllThreads returns all thread.
	GetAllThreads(ctx context.Context) ([]Thread, error)

	// GetThreadByID returns a thread with the given
	// thread ID.
	GetThreadByID(ctx context.Context, threadID int64) (Thread, error)

	// UpdateThread updates existing thread with the given
	// thread data.
	//
	// Updatethread do updates on all main attributes
	// except ID, CreateBy, and CreateTime. So, make sure to
	// use current values in the given data if do not want
	// to update some specific attributes.
	UpdateThread(ctx context.Context, thread Thread) error
}

type Thread struct {
	ID          int64
	Title       string
	GameID      int64
	GameNames   string // derived
	GameIcons   string // derived
	Description string
	ImageThread string
	Date        time.Time
	CreateTime  time.Time
	UpdateTime  time.Time
}
