package service

import (
	"context"

	"github.com/x-sports/internal/news"
)

func (s *service) CreateNews(ctx context.Context, reqNews news.News) (int64, error) {
	// validate field
	err := validateNews(reqNews)
	if err != nil {
		return 0, err
	}

	reqNews.CreateTime = s.timeNow()

	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return 0, err
	}

	newsID, err := pgStoreClient.CreateNews(ctx, reqNews)
	if err != nil {
		return 0, err
	}

	return newsID, nil
}

func (s *service) GetAllNews(ctx context.Context, gameID int64) ([]news.News, error) {
	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return nil, err
	}

	// get all news from postgre
	news, err := pgStoreClient.GetAllNews(ctx, gameID)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (s *service) UpdateNews(ctx context.Context, reqNews news.News) error {
	// validate field
	err := validateNews(reqNews)
	if err != nil {
		return err
	}

	// modify fields
	reqNews.UpdateTime = s.timeNow()

	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// updates news in pgstore
	err = pgStoreClient.UpdateNews(ctx, reqNews)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetNewsByID(ctx context.Context, newsID int64) (news.News, error) {
	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return news.News{}, err
	}

	// get news from pgstore
	result, err := pgStoreClient.GetNewsByID(ctx, newsID)
	if err != nil {
		return news.News{}, err
	}

	return result, nil
}

// validateNews validates fields of the given News
// whether its comply the predetermined rules.
func validateNews(reqNews news.News) error {
	if reqNews.Title == "" {
		return news.ErrInvalidTitle
	}

	if reqNews.GameID <= 0 {
		return news.ErrInvalidGameID
	}

	if reqNews.Description == "" {
		return news.ErrInvalidDescription
	}

	if reqNews.ImageNews == "" {
		return news.ErrInvalidImageNews
	}

	if reqNews.Date.IsZero() {
		return news.ErrInvalidDate
	}

	return nil
}
