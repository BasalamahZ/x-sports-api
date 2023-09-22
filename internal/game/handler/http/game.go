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
	"github.com/x-sports/internal/game"
)

type gameHandler struct {
	game  game.Service
	admin admin.Service
}

func (h *gameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Printf("[game HTTP][gameHandler] Failed to parse game ID. ID: %s. Err: %s\n", vars["id"], err.Error())
		helper.WriteErrorResponse(w, http.StatusBadRequest, []string{errInvalidGameID.Error()})
		return
	}

	switch r.Method {
	case http.MethodPatch:
		h.handleUpdateGame(w, r, gameID)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *gameHandler) handleUpdateGame(w http.ResponseWriter, r *http.Request, gameID int64) {
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
			log.Printf("[game HTTP][handleUpdateGame] Failed to update game. gameID: %d, Err: %s\n", gameID, err.Error())
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
		request := gameHTTP{}
		err = json.Unmarshal(body, &request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// check access token
		err = checkAccessToken(ctx, h.admin, token, "handleUpdateGame")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// get current curriculum
		current, err := h.game.GetGameByID(ctx, gameID)
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
				log.Printf("[game HTTP][handleUpdateGame] Internal error from GetGameByID. gameID: %d. Err: %s\n", gameID, err.Error())
			}

			errChan <- parsedErr
			return
		}

		// format HTTP request into service object
		reqGame, err := parseGameFromUpdateRequest(request, current)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		err = h.game.UpdateGame(ctx, reqGame)
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
				log.Printf("[Game HTTP][handleUpdateGame] Internal error from UpdateGame. Err: %s\n", err.Error())
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
			Data: gameID,
		})
	}
}

// parseGameFromUpdateRequest returns game.game from the
// given HTTP request object.
func parseGameFromUpdateRequest(gh gameHTTP, current game.Game) (game.Game, error) {
	result := current

	if gh.GameNames != nil {
		result.GameNames = *gh.GameNames
	}

	if gh.GameIcons != nil {
		result.GameIcons = *gh.GameIcons
	}

	return result, nil
}
