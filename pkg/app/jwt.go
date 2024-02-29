package app

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrMissingHeader = errors.New("the length of the `Authorization` header is zero")
	ErrInvalidToken  = errors.New("the `Authorization` header is not a valid JWT")
	ErrVertifyFailed = errors.New("the JWT failed to be verified")
	ErrInvalidClaims = errors.New("the JWT claims are invalid")
	ErrExpiredToken  = errors.New("the JWT has expired")
)

type Payload struct {
	UserID uint64
}

func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	}
}

func Parse(tokenString string, secret string) (*Payload, error) {
	token, err := jwt.Parse(tokenString, secretFunc(secret))
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &Payload{}, ErrInvalidClaims
	}

	if claims.VerifyExpiresAt(time.Now().Unix(), true) == false {
		return &Payload{}, ErrExpiredToken
	}

	if !token.Valid {
		return &Payload{}, ErrVertifyFailed
	}

	payloads := &Payload{}
	payloads.UserID = uint64(claims["user_id"].(float64))

	return payloads, nil
}

// Sign signs the payload with the specified secret.
// The token content.
// iss: （Issuer）签发者
// iat: （Issued At）签发时间，用Unix时间戳表示
// exp: （Expiration Time）过期时间，用Unix时间戳表示
// aud: （Audience）接收该JWT的一方
// sub: （Subject）该JWT的主题
// nbf: （Not Before）不要早于这个时间
// jti: （JWT ID）用于标识JWT的唯一ID
func Sign(ctx context.Context, payload map[string]interface{}, secret string, timeout time.Duration) (tokenString string, err error) {
	now := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["nbf"] = now
	claims["iat"] = now
	if timeout > 0 {
		claims["exp"] = now + int64(timeout)
	}

	for k, v := range payload {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the specified secret.
	tokenString, err = token.SignedString([]byte(secret))

	return
}
