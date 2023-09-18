package admin

import (
	"context"
	"time"
)

type Service interface {
	// LoginBasic checks the given email and password with
	// the actual data. It returns a token and the encapsulated
	// data if the login process is success.
	LoginBasic(ctx context.Context, email string, password string) (string, TokenData, error)

	// ValidateToken validates the given token and returns the
	// data encapsulated in the token if the given token is
	// valid.
	ValidateToken(ctx context.Context, token string) (TokenData, error)
}

type Admin struct {
	ID         int64
	Email      string
	Password   string
	CreateTime time.Time
	UpdateTime time.Time
}

// TokenData is the data that are encapsulated in a token.
type TokenData struct {
	AdminID int64
	Email   string
}
