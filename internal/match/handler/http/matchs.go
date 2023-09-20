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
	"github.com/x-sports/internal/match"
)

type matchsHandler struct {
	match match.Service
	admin admin.Service
}

func (h *matchsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetAllMatchs(w, r)
	case http.MethodPost:
		h.handleCreateMatch(w, r)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *matchsHandler) handleGetAllMatchs(w http.ResponseWriter, r *http.Request) {
	// add timeout to context
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Millisecond)
	defer cancel()

	var (
		err        error           // stores error in this handler
		source     string          // stores request source
		resBody    []byte          // stores response body to write
		statusCode = http.StatusOK // stores response status code
	)

	// write response
	defer func() {
		// error
		if err != nil {
			log.Printf("[Match HTTP][handleGetAllMatchs] Failed to get all matchs. Source: %s, Err: %s\n", source, err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan []match.Match, 1)
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
		err = checkAccessToken(ctx, h.admin, token, "handleGetAllMatchs")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// parsed filter
		gameID, err := parseGetMatchsFilter(r.URL.Query())
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		res, err := h.match.GetAllMatchs(ctx, gameID)
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
				log.Printf("[Match HTTP][handleGetAllGames] Internal error from GetAllMatchs. Err: %s\n", err.Error())
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
		// format each matchs
		matchs := make([]matchHTTP, 0)
		for _, r := range res {
			var m matchHTTP
			m, err = formatMatch(r)
			if err != nil {
				return
			}
			matchs = append(matchs, m)
		}

		// construct response data
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: matchs,
		})
	}
}

func parseGetMatchsFilter(request url.Values) (int64, error) {
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

func (h *matchsHandler) handleCreateMatch(w http.ResponseWriter, r *http.Request) {
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
			log.Printf("[Match HTTP][handleCreateMatch] Failed to create match. Err: %s\n", err.Error())
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
		request := matchHTTP{}
		err = json.Unmarshal(body, &request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// check access token
		err = checkAccessToken(ctx, h.admin, token, "handleCreateMatch")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// format HTTP request into service object
		reqMatch, err := parseMatchFromCreateRequest(request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		res, err := h.match.CreateMatch(ctx, reqMatch)
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
				log.Printf("[Match HTTP][handleCreateMatch] Internal error from CreateMatch. Err: %s\n", err.Error())
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
	case matchID := <-resChan:
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: matchID,
		})
	}
}

// parseMatchFromCreateRequest returns match.Match from the
// given HTTP request object.
func parseMatchFromCreateRequest(mh matchHTTP) (match.Match, error) {
	result := match.Match{
		Status: match.StatusUpcoming,
	}

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

	if mh.MatchLink != nil {
		result.MatchLink = *mh.MatchLink
	}

	return result, nil
}
