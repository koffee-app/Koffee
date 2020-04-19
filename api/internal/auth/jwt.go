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

// Information implements IUser
func (u *UserJWT) Information() (string, time.Time, uint32, int64) {
	return u.Email, u.LogedAt, u.UserID, u.ExpiresAt
}

type jwtService struct{}

// NewJWT token generator
func NewJWT() Token {
	return &jwtService{}
}

const expiresHours = 24

// ExpiresAt Returns expires at number
func ExpiresAt() int64 {
	return time.Now().Add(time.Hour * expiresHours).Unix()
}

// newUser Returns a new UserJWT
func (j *jwtService) newUser(email string, id uint32) *UserJWT {
	user := &UserJWT{Email: email, UserID: id, LogedAt: time.Now()}

	return user
}

// FormatSpecifics returns Bearer + token
func (j *jwtService) FormatSpecifics(t string) string {
	return fmt.Sprintf("Bearer %s", t)
}

// GenerateToken Returns the string of the generated token.
func (j *jwtService) GenerateToken(email string, id uint32, duration uint64) (string, error) {
	secret := config.JWTKey()
	user := j.newUser(email, id)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	user.ExpiresAt = time.Now().Add(time.Minute * time.Duration(duration)).Unix()
	return claims.SignedString([]byte(secret))
}

// TokenValid str
func (j *jwtService) TokenValid(r *http.Request) (IUser, error) {
	token := j.ParseToken(r)
	return j.VerifyToken(token)
}

// VerifyToken Verifies a token string
func (j *jwtService) VerifyToken(tokenStr string) (IUser, error) {
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
func (j *jwtService) ParseToken(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	splittedBearer := strings.Split(bearer, " ")
	if len(splittedBearer) < 2 {
		return ""
	}

	return splittedBearer[1]
}
