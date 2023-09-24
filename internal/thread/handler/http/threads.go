package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/x-sports/global/helper"
	"github.com/x-sports/internal/admin"
	"github.com/x-sports/internal/thread"
)

type threadsHandler struct {
	thread thread.Service
	admin  admin.Service
}

func (h *threadsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetAllThreads(w, r)
	case http.MethodPost:
		h.handleCreateThread(w, r)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *threadsHandler) handleGetAllThreads(w http.ResponseWriter, r *http.Request) {
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
			log.Printf("[Thread HTTP][handleGetAllThreads] Failed to get all thread. Err: %s\n", err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan []thread.Thread, 1)
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
		err = checkAccessToken(ctx, h.admin, token, "handleGetAllThreads")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		res, err := h.thread.GetAllThreads(ctx)
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
				log.Printf("[Thread HTTP][handleGetAllGames] Internal error from GetAllThreads. Err: %s\n", err.Error())
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
		// format each thread
		thread := make([]threadHTTP, 0)
		for _, r := range res {
			var t threadHTTP
			t, err = formatThread(r)
			if err != nil {
				return
			}
			thread = append(thread, t)
		}

		// construct response data
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: thread,
		})
	}
}

func (h *threadsHandler) handleCreateThread(w http.ResponseWriter, r *http.Request) {
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
			log.Printf("[Thread HTTP][handleCreateThread] Failed to create thread. Err: %s\n", err.Error())
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
		request := threadHTTP{}
		err = json.Unmarshal(body, &request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// check access token
		err = checkAccessToken(ctx, h.admin, token, "handleCreateThread")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// format HTTP request into service object
		reqThread, err := parseThreadFromCreateRequest(request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		res, err := h.thread.CreateThread(ctx, reqThread)
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
				log.Printf("[Thread HTTP][handleCreateThread] Internal error from CreateThread. Err: %s\n", err.Error())
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
	case threadID := <-resChan:
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: threadID,
		})
	}
}

// parseThreadFromCreateRequest returns thread.Thread from the
// given HTTP request object.
func parseThreadFromCreateRequest(nh threadHTTP) (thread.Thread, error) {
	result := thread.Thread{}

	if nh.Title != nil {
		result.Title = *nh.Title
	}

	if nh.GameID != nil {
		result.GameID = *nh.GameID
	}

	if nh.Description != nil {
		result.Description = *nh.Description
	}

	if nh.ImageThread != nil {
		result.ImageThread = *nh.ImageThread
	}

	if nh.Date != nil && *nh.Date != "" {
		date, err := time.Parse(dateFormat, *nh.Date)
		if err != nil {
			return thread.Thread{}, errInvalidTimeFormat
		}
		result.Date = date
	}

	return result, nil
}
