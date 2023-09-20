package http

import (
	"context"
	"log"

	"github.com/x-sports/internal/admin"
)

// checkAccessToken checks the given access token whether it
// is valid or not.
func checkAccessToken(ctx context.Context, svc admin.Service, token, name string) error {
	_, err := svc.ValidateToken(ctx, token)
	if err != nil {
		log.Printf("[Game HTTP][%s] Unauthorized error from ValidateToken. Err: %s\n", name, err.Error())
		return errUnauthorizedAccess
	}

	return nil
}
