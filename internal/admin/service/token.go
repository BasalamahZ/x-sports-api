package service

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/x-sports/internal/admin"
)

// jwtClaimss is the claims encapsulated in JWT-generated token.
type jwtClaims struct {
	AdminID int64  `json:"admin_id"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}

// parseTokenData parse token data from jwt claims.
func (jwtc jwtClaims) parseTokenData() admin.TokenData {
	return admin.TokenData{
		AdminID: jwtc.AdminID,
		Email:   jwtc.Email,
	}
}

// formatTokenData format token data into jwt claims.
func formatTokenData(data admin.TokenData) jwtClaims {
	return jwtClaims{
		AdminID: data.AdminID,
		Email:   data.Email,
	}
}

func (s *service) LoginBasic(ctx context.Context, email string, password string) (string, admin.TokenData, error) {
	// validate the given values
	if email == "" {
		return "", admin.TokenData{}, admin.ErrInvalidEmail
	}
	if password == "" {
		return "", admin.TokenData{}, admin.ErrInvalidPassword
	}

	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return "", admin.TokenData{}, err
	}

	// get user current data
	current, err := pgStoreClient.GetUserByEmail(ctx, email)
	if err != nil {
		return "", admin.TokenData{}, admin.ErrInvalidEmail
	}

	// check password
	if password != current.Password {
		return "", admin.TokenData{}, admin.ErrInvalidPassword
	}

	// generate token
	tokenData := admin.TokenData{
		AdminID: current.ID,
		Email:   current.Email,
	}
	token, err := s.generateToken(ctx, tokenData)
	if err != nil {
		return "", admin.TokenData{}, err
	}

	return token, tokenData, nil
}

// generateToken returns a new token that encapsulates the
// given token data with some additional information:
//   - token expiration time
//
// Token is generated using JWT HS256.
func (s *service) generateToken(ctx context.Context, data admin.TokenData) (string, error) {
	claims := formatTokenData(data)

	// add expirations time
	expiresAt := s.timeNow().Add(s.config.TokenExpiration)
	claims.ExpiresAt = jwt.NewNumericDate(expiresAt)

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign token with secret key
	signedToken, err := token.SignedString([]byte(s.config.TokenSecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// TODO: check for expired token error from internal JWT library
func (s *service) ValidateToken(ctx context.Context, token string) (admin.TokenData, error) {
	if token == "" {
		return admin.TokenData{}, admin.ErrInvalidToken
	}

	// get jwt token object
	jwtToken, err := jwt.ParseWithClaims(token, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.config.TokenSecretKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return admin.TokenData{}, admin.ErrInvalidToken
		}
		return admin.TokenData{}, err
	}

	// check whether token is valid or not (from expirations time)
	if !jwtToken.Valid {
		return admin.TokenData{}, admin.ErrExpiredToken
	}

	// parse jwt claims
	claims, ok := jwtToken.Claims.(*jwtClaims)
	if !ok {
		return admin.TokenData{}, admin.ErrInvalidToken
	}

	return claims.parseTokenData(), nil
}
