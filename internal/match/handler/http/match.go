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
	"github.com/x-sports/internal/match"
)

type matchHandler struct {
	match match.Service
	admin admin.Service
}

func (h *matchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	matchID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Printf("[Match HTTP][matchHandler] Failed to parse match ID. ID: %s. Err: %s\n", vars["id"], err.Error())
		helper.WriteErrorResponse(w, http.StatusBadRequest, []string{errInvalidMatchID.Error()})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetMatchByID(w, r, matchID)
	case http.MethodPatch:
		h.handleUpdateMatch(w, r, matchID)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *matchHandler) handleGetMatchByID(w http.ResponseWriter, r *http.Request, matchID int64) {
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
			log.Printf("[Match HTTP][handleGetMatchByID] Failed to get match by ID. matchID: %d, Err: %s\n", matchID, err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan match.Match, 1)
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
		err = checkAccessToken(ctx, h.admin, token, "handleGetMatchByID")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// TODO: add authorization flow with roles

		res, err := h.match.GetMatchByID(ctx, matchID)
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
				log.Printf("[Match HTTP][handleGetMatchByID] Internal error from GetMatchByID. matchID: %d. Err: %s\n", matchID, err.Error())
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
		var m matchHTTP
		m, err = formatMatch(res)
		if err != nil {
			return
		}
		// construct response data
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: m,
		})
	}
}

func (h *matchHandler) handleUpdateMatch(w http.ResponseWriter, r *http.Request, matchID int64) {
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
			log.Printf("[Match HTTP][handleUpdateMatch] Failed to update match. matchID: %d, Err: %s\n", matchID, err.Error())
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
		request := matchHTTP{}
		err = json.Unmarshal(body, &request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// check access token
		err = checkAccessToken(ctx, h.admin, token, "handleUpdateMatch")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// get current curriculum
		current, err := h.match.GetMatchByID(ctx, matchID)
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
				log.Printf("[Match HTTP][handleUpdateMatch] Internal error from GetMatchByID. matchID: %d. Err: %s\n", matchID, err.Error())
			}

			errChan <- parsedErr
			return
		}

		// format HTTP request into service object
		reqMatch, err := parseMatchFromUpdateRequest(request, current)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		err = h.match.UpdateMatch(ctx, reqMatch)
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
				log.Printf("[Match HTTP][handleUpdateMatch] Internal error from UpdateMatch. Err: %s\n", err.Error())
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
			Data: matchID,
		})
	}
}

// parseMatchFromUpdateRequest returns match.Match from the
// given HTTP request object.
func parseMatchFromUpdateRequest(mh matchHTTP, current match.Match) (match.Match, error) {
	result := current

	if mh.TournamentNames != nil {
		result.TournamentNames = *mh.TournamentNames
	}

	if mh.GameID != nil {
		result.GameID = *mh.GameID
	}

	if mh.TeamAID != nil {
		result.TeamAID = *mh.TeamAID
	}

	if mh.TeamBID != nil {
		result.TeamBID = *mh.TeamBID
	}

	if mh.TeamAOdds != nil {
		result.TeamAOdds = *mh.TeamAOdds
	}

	if mh.TeamBOdds != nil {
		result.TeamBOdds = *mh.TeamBOdds
	}

	if mh.Date != nil && *mh.Date != "" {
		date, err := time.Parse(dateFormat, *mh.Date)
		if err != nil {
			return match.Match{}, errInvalidTimeFormat
		}
		result.Date = date
	}

	if mh.Status != nil {
		status, err := parseStatus(*mh.Status)
		if err != nil {
			return match.Match{}, err
		}
		result.Status = status
	}

	if mh.MatchLink != nil {
		result.MatchLink = *mh.MatchLink
	}

	if mh.Winner != nil {
		result.Winner = *mh.Winner
	}

	return result, nil
}
