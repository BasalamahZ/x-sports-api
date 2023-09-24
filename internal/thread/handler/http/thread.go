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
	"github.com/x-sports/internal/thread"
)

type threadHandler struct {
	thread thread.Service
	admin  admin.Service
}

func (h *threadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Printf("[Thread HTTP][threadHandler] Failed to parse thread ID. ID: %s. Err: %s\n", vars["id"], err.Error())
		helper.WriteErrorResponse(w, http.StatusBadRequest, []string{errInvalidThreadID.Error()})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetThreadByID(w, r, threadID)
	case http.MethodPatch:
		h.handleUpdateThread(w, r, threadID)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *threadHandler) handleGetThreadByID(w http.ResponseWriter, r *http.Request, threadID int64) {
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
			log.Printf("[Thread HTTP][handleGetThreadByID] Failed to get thread by ID. threadID: %d, Err: %s\n", threadID, err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan thread.Thread, 1)
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
		err = checkAccessToken(ctx, h.admin, token, "handleGetThreadByID")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// TODO: add authorization flow with roles

		res, err := h.thread.GetThreadByID(ctx, threadID)
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
				log.Printf("[Thread HTTP][handleGetThreadByID] Internal error from GetThreadByID. threadID: %d. Err: %s\n", threadID, err.Error())
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
		var t threadHTTP
		t, err = formatThread(res)
		if err != nil {
			return
		}
		// construct response data
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: t,
		})
	}
}

func (h *threadHandler) handleUpdateThread(w http.ResponseWriter, r *http.Request, threadID int64) {
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
			log.Printf("[Thread HTTP][handleUpdateThread] Failed to update thread. threadID: %d, Err: %s\n", threadID, err.Error())
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
		request := threadHTTP{}
		err = json.Unmarshal(body, &request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// check access token
		err = checkAccessToken(ctx, h.admin, token, "handleUpdateThread")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// get current thread
		current, err := h.thread.GetThreadByID(ctx, threadID)
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
				log.Printf("[Thread HTTP][handleUpdateThread] Internal error from GetThreadByID. threadID: %d. Err: %s\n", threadID, err.Error())
			}

			errChan <- parsedErr
			return
		}

		// format HTTP request into service object
		reqThread, err := parseThreadFromUpdateRequest(request, current)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		err = h.thread.UpdateThread(ctx, reqThread)
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
				log.Printf("[Thread HTTP][handleUpdateThread] Internal error from UpdateThread. Err: %s\n", err.Error())
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
			Data: threadID,
		})
	}
}

// parseThreadFromUpdateRequest returns Thread.Thread from the
// given HTTP request object.
func parseThreadFromUpdateRequest(nh threadHTTP, current thread.Thread) (thread.Thread, error) {
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
