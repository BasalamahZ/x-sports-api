package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/x-sports/global/helper"
	"github.com/x-sports/internal/admin"
	"github.com/x-sports/internal/news"
)

type newsHandler struct {
	news  news.Service
	admin admin.Service
}

func (h *newsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	newsID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Printf("[News HTTP][newsHandler] Failed to parse news ID. ID: %s. Err: %s\n", vars["id"], err.Error())
		helper.WriteErrorResponse(w, http.StatusBadRequest, []string{errInvalidNewsID.Error()})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetNewsByID(w, r, newsID)
	case http.MethodPatch:
		h.handleUpdateNews(w, r, newsID)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *newsHandler) handleGetNewsByID(w http.ResponseWriter, r *http.Request, newsID int64) {
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
			log.Printf("[News HTTP][handleGetNewsByID] Failed to get news by ID. newsID: %d, Err: %s\n", newsID, err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan news.News, 1)
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
		err = checkAccessToken(ctx, h.admin, token, "handleGetNewsByID")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// TODO: add authorization flow with roles

		res, err := h.news.GetNewsByID(ctx, newsID)
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
				log.Printf("[News HTTP][handleGetNewsByID] Internal error from GetNewsByID. newsID: %d. Err: %s\n", newsID, err.Error())
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
		// format curriculum
		var n newsHTTP
		n, err = formatNews(res)
		if err != nil {
			return
		}
		// construct response data
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: n,
		})
	}
}

func (h *newsHandler) handleUpdateNews(w http.ResponseWriter, r *http.Request, newsID int64) {
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
			log.Printf("[News HTTP][handleUpdateNews] Failed to update news. newsID: %d, Err: %s\n", newsID, err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan struct{}, 1)
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
		err = checkAccessToken(ctx, h.admin, token, "handleUpdateNews")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// get current curriculum
		current, err := h.news.GetNewsByID(ctx, newsID)
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
				log.Printf("[News HTTP][handleUpdateNews] Internal error from GetNewsByID. newsID: %d. Err: %s\n", newsID, err.Error())
			}

			errChan <- parsedErr
			return
		}

		// format HTTP request into service object
		reqNews, err := parseNewsFromUpdateRequest(request, current)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		err = h.news.UpdateNews(ctx, reqNews)
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
				log.Printf("[News HTTP][handleUpdateNews] Internal error from UpdateNews. Err: %s\n", err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- struct{}{}
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case <-resChan:
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: newsID,
		})
	}
}

// parseNewsFromUpdateRequest returns news.News from the
// given HTTP request object.
func parseNewsFromUpdateRequest(nh newsHTTP, current news.News) (news.News, error) {
	result := current

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
