package auth

import (
	"fmt"
	"time"
	"tsumiki/env"

	"github.com/bwmarrin/snowflake"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// AccessTokenLiveTime  = 15 * time.Minute
	AccessTokenLiveTime  = 150000 * time.Minute // 動作確認用
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

func ValidateRefreshToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return env.JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func GenerateTokenPair(userID int) (TokenPair, error) {
	sessionID := snowflakeNode.Generate().String()
	now := time.Now()

	accessClaims := CustomClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenLiveTime)),
			IssuedAt:  jwt.NewNumericDate(now),
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
			ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenLiveTime)),
			IssuedAt:  jwt.NewNumericDate(now),
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
