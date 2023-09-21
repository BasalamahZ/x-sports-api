package http

import "github.com/x-sports/internal/news"

// formatNews formats the given news into the
// respective HTTP-format object.
func formatNews(n news.News) (newsHTTP, error) {
	date := n.Date.Format(dateFormat)

	return newsHTTP{
		ID:          &n.ID,
		Title:       &n.Title,
		GameID:      &n.GameID,
		GameNames:   &n.GameNames,
		GameIcons:   &n.GameIcons,
		Description: &n.Description,
		ImageNews:   &n.ImageNews,
		Date:        &date,
	}, nil
}
