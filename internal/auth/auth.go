package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// TokenResponse represents the response when issuing a token
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// JWTManager handles JWT operations
type JWTManager struct {
	secretKey       string
	refreshSecretKey string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey, refreshSecretKey string, accessTokenTTL, refreshTokenTTL time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:        secretKey,
		refreshSecretKey: refreshSecretKey,
		accessTokenTTL:   accessTokenTTL,
		refreshTokenTTL:  refreshTokenTTL,
	}
}

// GenerateTokens generates both access and refresh tokens
func (jm *JWTManager) GenerateTokens(userID int, username, email string) (TokenResponse, error) {
	now := time.Now()
	expiresAt := now.Add(jm.accessTokenTTL)

	claims := Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// Create access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(jm.secretKey))
	if err != nil {
		return TokenResponse{}, err
	}

	// Create refresh token
	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(jm.refreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(jm.refreshSecretKey))
	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(jm.accessTokenTTL.Seconds()),
	}, nil
}

// ValidateToken validates the access token
func (jm *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jm.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token from a refresh token
func (jm *JWTManager) RefreshAccessToken(refreshTokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(refreshTokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jm.refreshSecretKey), nil
	})

	if err != nil {
		return "", err
	}

	_, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	// Create new access token with same claims but new expiration
	now := time.Now()
	newClaims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jm.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return newAccessToken.SignedString([]byte(jm.secretKey))
}
