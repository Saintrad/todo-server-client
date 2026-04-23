package auth

import (
	"time"

	"github.com/Saintrad/todo-server-client/internal/richerror"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
    secretKey     string
    tokenDuration time.Duration
}

func NewJWTManager(secret string, duration time.Duration) *JWTManager {
    return &JWTManager{
        secretKey:     secret,
        tokenDuration: duration,
    }
}

func (m *JWTManager) GenerateToken(userID int) (string, *richerror.AppError) {
    claims := jwt.MapClaims{
        "uid": userID,
        "exp": time.Now().Add(m.tokenDuration).Unix(),
        "iat": time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", richerror.Internal("token generation failed", err)
	}

	return signed, nil
}

func (m *JWTManager) VerifyToken(tokenStr string) (*jwt.Token, *richerror.AppError) {
    token, err :=jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
        return []byte(m.secretKey), nil
    })

	if err != nil {
		return nil, richerror.Internal("token verification failed", err)
	}

	return token, nil
}
