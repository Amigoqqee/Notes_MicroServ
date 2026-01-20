package jwtmanager

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	ACCESS_TOKEN  = "accessToken"
	REFRESH_TOKEN = "refreshToken"
)

type JWTConfig struct {
	SecretKey              string
	AccessTokenExpiration  int
	RefreshTokenExpiration int
}

type JWTManager struct {
	config JWTConfig
}

func NewJWTManager(config JWTConfig) *JWTManager {
	return &JWTManager{
		config: config,
	}
}

func (s *JWTManager) GenerateTokens(id int) (access, refresh string, err error) {
	accessTokenString, err := s.generateToken(id, ACCESS_TOKEN, s.config.AccessTokenExpiration)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", ErrTokenGeneration, err)
	}

	refreshTokenString, err := s.generateToken(id, REFRESH_TOKEN, s.config.RefreshTokenExpiration)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", ErrTokenGeneration, err)
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *JWTManager) ValidateAccessToken(tokenString string) (int, error) {
	return s.validateToken(tokenString, ACCESS_TOKEN)
}

func (s *JWTManager) ValidateRefreshToken(tokenString string) (int, error) {
	return s.validateToken(tokenString, REFRESH_TOKEN)
}

func (s *JWTManager) generateToken(id int, tokenType string, expirationHours int) (string, error) {
	now := time.Now()
	expiration := now.Add(time.Hour * time.Duration(expirationHours))

	claims := jwt.MapClaims{
		"id":   id,
		"type": tokenType,
		"iat":  now.Unix(),
		"exp":  expiration.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrInvalidSignature, err)
	}

	return tokenString, nil
}

func (s *JWTManager) validateToken(tokenString, tokenType string) (int, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: %v", ErrInvalidSignature, token.Header["alg"])
		}
		return []byte(s.config.SecretKey), nil
	}

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return 0, fmt.Errorf("%s: %w", ErrTokenExpired, err)
			}
		}
		return 0, fmt.Errorf("%s: %w", ErrInvalidToken, err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if tokenType != "" {
			if claimType, exists := claims["type"].(string); !exists || claimType != tokenType {
				return 0, fmt.Errorf("%s: ожидается %s, получен %s", ErrInvalidTokenType, tokenType, claims["type"])
			}
		}
		idValue, exists := claims["id"].(float64)
		if !exists {
			return 0, fmt.Errorf("%s", ErrMissingUserID)
		}

		return int(idValue), nil
	}

	return 0, fmt.Errorf("%s", ErrInvalidToken)
}
