package service

import (
	"context"

	"github.com/x-sports/internal/thread"
)

func (s *service) CreateThread(ctx context.Context, reqThread thread.Thread) (int64, error) {
	// validate field
	err := validateThread(reqThread)
	if err != nil {
		return 0, err
	}

	reqThread.CreateTime = s.timeNow()

	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return 0, err
	}

	threadID, err := pgStoreClient.CreateThread(ctx, reqThread)
	if err != nil {
		return 0, err
	}

	return threadID, nil
}

func (s *service) GetAllThreads(ctx context.Context) ([]thread.Thread, error) {
	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return nil, err
	}

	// get all thread from postgre
	thread, err := pgStoreClient.GetAllThreads(ctx)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (s *service) UpdateThread(ctx context.Context, reqThread thread.Thread) error {
	// validate field
	err := validateThread(reqThread)
	if err != nil {
		return err
	}

	// modify fields
	reqThread.UpdateTime = s.timeNow()

	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// updates thread in pgstore
	err = pgStoreClient.UpdateThread(ctx, reqThread)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetThreadByID(ctx context.Context, threadID int64) (thread.Thread, error) {
	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return thread.Thread{}, err
	}

	// get thread from pgstore
	result, err := pgStoreClient.GetThreadByID(ctx, threadID)
	if err != nil {
		return thread.Thread{}, err
	}

	return result, nil
}

// validateThread validates fields of the given Thread
// whether its comply the predetermined rules.
func validateThread(reqThread thread.Thread) error {
	if reqThread.Title == "" {
		return thread.ErrInvalidTitle
	}

	if reqThread.GameID <= 0 {
		return thread.ErrInvalidGameID
	}

	if reqThread.Description == "" {
		return thread.ErrInvalidDescription
	}

	if reqThread.ImageThread == "" {
		return thread.ErrInvalidImageThread
	}

	if reqThread.Date.IsZero() {
		return thread.ErrInvalidDate
	}

	return nil
}
