package middleware

import (
	"fmt"
	"time"
	"tsumiki/env"

	"github.com/bwmarrin/snowflake"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenLiveTime  = 15 * time.Minute
	RefreshTokenLiveTime = 30 * 24 * time.Hour
)

var snowflakeNode *snowflake.Node

func init() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(fmt.Sprintf("snowflake node init failed: %v", err))
	}
	snowflakeNode = node
}

type CustomClaims struct {
	UserID    int    `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	UserID       int
	SessionID    string
}

func GenerateTokenPair(userID int) (TokenPair, error) {
	sessionID := snowflakeNode.Generate().String()

	accessClaims := CustomClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenLiveTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(env.JwtSecret)
	if err != nil {
		return TokenPair{}, err
	}

	refreshClaims := CustomClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenLiveTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(env.JwtSecret)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		UserID:       userID,
		SessionID:    sessionID,
	}, nil
}
