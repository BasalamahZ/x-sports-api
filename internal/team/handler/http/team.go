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
	"github.com/x-sports/internal/team"
)

type teamHandler struct {
	team  team.Service
	admin admin.Service
}

func (h *teamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Printf("[Team HTTP][teamHandler] Failed to parse team ID. ID: %s. Err: %s\n", vars["id"], err.Error())
		helper.WriteErrorResponse(w, http.StatusBadRequest, []string{errInvalidTeamID.Error()})
		return
	}

	switch r.Method {
	case http.MethodPatch:
		h.handleUpdateTeam(w, r, teamID)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *teamHandler) handleUpdateTeam(w http.ResponseWriter, r *http.Request, teamID int64) {
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
			log.Printf("[Team HTTP][handleUpdateTeam] Failed to update team. teamID: %d, Err: %s\n", teamID, err.Error())
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
		request := teamHTTP{}
		err = json.Unmarshal(body, &request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// check access token
		err = checkAccessToken(ctx, h.admin, token, "handleUpdateTeam")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// get current curriculum
		current, err := h.team.GetTeamByID(ctx, teamID)
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
				log.Printf("[Team HTTP][handleUpdateTeam] Internal error from GetTeamByID. teamID: %d. Err: %s\n", teamID, err.Error())
			}

			errChan <- parsedErr
			return
		}

		// format HTTP request into service object
		reqTeam, err := parseTeamFromUpdateRequest(request, current)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		err = h.team.UpdateTeam(ctx, reqTeam)
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
				log.Printf("[Team HTTP][handleUpdateTeam] Internal error from UpdateTeam. Err: %s\n", err.Error())
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
			Data: teamID,
		})
	}
}

// parseTeamFromUpdateRequest returns team.Team from the
// given HTTP request object.
func parseTeamFromUpdateRequest(th teamHTTP, current team.Team) (team.Team, error) {
	result := current

	if th.TeamNames != nil {
		result.TeamNames = *th.TeamNames
	}

	if th.TeamIcons != nil {
		result.TeamIcons = *th.TeamIcons
	}

	if th.GameID != nil {
		result.GameID = *th.GameID
	}

	return result, nil
}
