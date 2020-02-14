package auth

import (
	"fmt"
	"koffee/internal/config"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// UserJWT stores current JWT Payload
type UserJWT struct {
	Email   string
	LogedAt time.Time
	UserID  uint32

	jwt.StandardClaims
}

const expiresHours = 24

// ExpiresAt Returns expires at number
func ExpiresAt() int64 {
	return time.Now().Add(time.Hour * expiresHours).Unix()
}

// NewUserJWT Returns a new UserJWT
func NewUserJWT(email string, id uint32) *UserJWT {
	user := &UserJWT{Email: email, UserID: id, LogedAt: time.Now()}
	user.ExpiresAt = time.Now().Add(time.Hour * expiresHours).Unix()
	return user
}

// GenerateTokenJWT Returns the string of the generated token.
func GenerateTokenJWT(email string, id uint32) (string, error) {
	secret := config.JWTKey()
	user := NewUserJWT(email, id)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, user)

	return claims.SignedString([]byte(secret))
}

// TokenValid str
func TokenValid(r *http.Request) (*UserJWT, error) {
	tokenStr := ParseToken(r)

	if tokenStr == "" {
		return nil, fmt.Errorf("Invalid token, empty or invalid format in Authorization header")
	}

	tk := &UserJWT{}

	token, err := jwt.ParseWithClaims(tokenStr, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTKey()), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	return tk, nil
}

// ParseToken Gets token from the request
func ParseToken(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	splittedBearer := strings.Split(bearer, " ")
	if len(splittedBearer) < 2 {
		return ""
	}

	return splittedBearer[1]
}
