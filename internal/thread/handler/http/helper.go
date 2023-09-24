package http

import "github.com/x-sports/internal/thread"

// formatThread formats the given thread into the
// respective HTTP-format object.
func formatThread(t thread.Thread) (threadHTTP, error) {
	date := t.Date.Format(dateFormat)

	return threadHTTP{
		ID:          &t.ID,
		Title:       &t.Title,
		GameID:      &t.GameID,
		GameNames:   &t.GameNames,
		GameIcons:   &t.GameIcons,
		Description: &t.Description,
		ImageThread: &t.ImageThread,
		Date:        &date,
	}, nil
}
