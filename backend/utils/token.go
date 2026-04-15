package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPayload struct {
	UserID string `json:"_id"`
	Role   string `json:"role"`
}

type Claims struct {
	UserID string `json:"_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret    string
	expiresIn time.Duration
}

func NewTokenService() *TokenService {

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET is not defined in .env")
	}

	exp := os.Getenv("JWT_EXPIRES_IN")

	var duration time.Duration
	if exp == "" {
		duration = time.Hour
	} else {
		parsed, err := time.ParseDuration(exp)
		if err != nil {
			duration = time.Hour
		} else {
			duration = parsed
		}
	}

	return &TokenService{
		secret:    secret,
		expiresIn: duration,
	}
}

func (t *TokenService) GenerateToken(payload TokenPayload) (string, error) {

	claims := Claims{
		UserID: payload.UserID,
		Role:   payload.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(t.secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (t *TokenService) VerifyToken(tokenStr string) (*TokenPayload, error) {

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(t.secret), nil
		},
	)

	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return &TokenPayload{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}

func ExtractTokenFromHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}
