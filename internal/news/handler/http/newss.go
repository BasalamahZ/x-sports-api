package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/x-sports/global/helper"
	"github.com/x-sports/internal/admin"
	"github.com/x-sports/internal/news"
)

type newssHandler struct {
	news  news.Service
	admin admin.Service
}

func (h *newssHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetAllNews(w, r)
	case http.MethodPost:
		h.handleCreateNews(w, r)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *newssHandler) handleGetAllNews(w http.ResponseWriter, r *http.Request) {
	// add timeout to context
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Millisecond)
	defer cancel()

	var (
		err        error           // stores error in this handler
		resBody    []byte          // stores response body to write
		statusCode = http.StatusOK // stores response status code
	)

	// write response
	defer func() {
		// error
		if err != nil {
			log.Printf("[News HTTP][handleGetAllNews] Failed to get all news. Err: %s\n", err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan []news.News, 1)
	errChan := make(chan error, 1)

	go func() {
		// get token from header
		token, err := helper.GetBearerTokenFromHeader(r)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errInvalidToken
			return
		}

		// check access token
		err = checkAccessToken(ctx, h.admin, token, "handleGetAllNews")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// parsed filter
		gameID, err := parseGetNewsFilter(r.URL.Query())
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		res, err := h.news.GetAllNews(ctx, gameID)
		if err != nil {
			// determine error and status code, by default its internal error
			parsedErr := errInternalServer
			statusCode = http.StatusInternalServerError
			if v, ok := mapHTTPError[err]; ok {
				parsedErr = v
				statusCode = http.StatusBadRequest
			}

			// log the actual error if its internal error
			if statusCode == http.StatusInternalServerError {
				log.Printf("[News HTTP][handleGetAllGames] Internal error from GetAllNews. Err: %s\n", err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- res
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case res := <-resChan:
		// format each news
		news := make([]newsHTTP, 0)
		for _, r := range res {
			var m newsHTTP
			m, err = formatNews(r)
			if err != nil {
				return
			}
			news = append(news, m)
		}

		// construct response data
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: news,
		})
	}
}

func parseGetNewsFilter(request url.Values) (int64, error) {
	var gameID int64
	if gameIDStr := request.Get("game_id"); gameIDStr != "" {
		intSGameID, err := strconv.ParseInt(gameIDStr, 10, 64)
		if err != nil {
			return 0, errInvalidGameID
		}
		gameID = intSGameID
	}

	return gameID, nil
}

func (h *newssHandler) handleCreateNews(w http.ResponseWriter, r *http.Request) {
	// add timeout to context
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Millisecond)
	defer cancel()

	var (
		err        error           // stores error in this handler
		resBody    []byte          // stores response body to write
		statusCode = http.StatusOK // stores response status code
	)

	// write response
	defer func() {
		// error
		if err != nil {
			log.Printf("[News HTTP][handleCreateNews] Failed to create news. Err: %s\n", err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan int64, 1)
	errChan := make(chan error, 1)

	go func() {
		// read request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// get token from header
		token, err := helper.GetBearerTokenFromHeader(r)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errInvalidToken
			return
		}

		// unmarshall body
		request := newsHTTP{}
		err = json.Unmarshal(body, &request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// check access token
		err = checkAccessToken(ctx, h.admin, token, "handleCreateNews")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// format HTTP request into service object
		reqNews, err := parseNewsFromCreateRequest(request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		res, err := h.news.CreateNews(ctx, reqNews)
		if err != nil {
			// determine error and status code, by default its internal error
			parsedErr := errInternalServer
			statusCode = http.StatusInternalServerError
			if v, ok := mapHTTPError[err]; ok {
				parsedErr = v
				statusCode = http.StatusBadRequest
			}

			// log the actual error if its internal error
			if statusCode == http.StatusInternalServerError {
				log.Printf("[News HTTP][handleCreateNews] Internal error from CreateNews. Err: %s\n", err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- res
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case newsID := <-resChan:
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: newsID,
		})
	}
}

// parseNewsFromCreateRequest returns news.News from the
// given HTTP request object.
func parseNewsFromCreateRequest(nh newsHTTP) (news.News, error) {
	result := news.News{}

	if nh.Title != nil {
		result.Title = *nh.Title
	}

	if nh.GameID != nil {
		result.GameID = *nh.GameID
	}

	if nh.Description != nil {
		result.Description = *nh.Description
	}

	if nh.ImageNews != nil {
		result.ImageNews = *nh.ImageNews
	}

	if nh.Date != nil && *nh.Date != "" {
		date, err := time.Parse(dateFormat, *nh.Date)
		if err != nil {
			return news.News{}, errInvalidTimeFormat
		}
		result.Date = date
	}

	return result, nil
}
