package auth

import (
	// MIN GO ~1.13

	"fmt"
	"koffee/internal/config"
	"net/http"
	"strconv"
	"time"

	"github.com/o1egl/paseto"
)

//
// type Token interface {
// 	// Generates a token
// 	GenerateToken(string, uint32) (string, error)
// 	// Checks if the token included in the request is valid
// 	TokenValid(*http.Request) (IUser, error)
// 	// Verifies token
// 	VerifyToken(tokenStr string) (IUser, error)
// 	// Gets token from request
// 	ParseToken(r *http.Request) string

// 	FormatSpecifics(string) string
// }

type userPaseto struct {
	Email     string
	LogedAt   time.Time
	UserID    uint32
	ExpiresAt int64
}

// NewPaseto .
func NewPaseto() Token {
	return &pasetoService{}
}

func (u *userPaseto) Information() (string, time.Time, uint32, int64) {
	return u.Email, u.LogedAt, u.UserID, u.ExpiresAt
}

type pasetoService struct{}

func (p *pasetoService) GenerateToken(email string, id uint32, duration uint64) (string, error) {
	idStr := fmt.Sprint(id)

	privateKey := config.PrivateKeyParsed()
	token := paseto.JSONToken{
		Expiration: time.Now().Add(time.Duration(duration) * time.Minute),
	}

	token.Set("email", email)
	token.Set("id", idStr)
	tokenStr, err := paseto.NewV2().Sign(privateKey, token, "")
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (p *pasetoService) TokenValid(r *http.Request) (IUser, error) {
	s := r.Header.Get("Authorization")
	val, err := r.Cookie("token")
	if err != nil {
		return nil, fmt.Errorf("Error retrieving cookies, %v", err)
	}
	if s != val.Value {
		return nil, fmt.Errorf("Error, cookies token are different from Authorization token: %s", val.Value)
	}
	return p.VerifyToken(s)
}

func (p *pasetoService) VerifyToken(tokenStr string) (IUser, error) {
	// symmetricKey := []byte(config.JWTKey())
	var token paseto.JSONToken
	var newFooter string
	// err := paseto.NewV2().Decrypt(tokenStr, []byte(config.JWTKey()), &token, &newFooter)
	publicKey := config.PublicKeyParsed()
	err := paseto.NewV2().Verify(tokenStr, publicKey, &token, &newFooter)
	if err != nil {
		return nil, err
	}

	user := &userPaseto{}
	user.Email = token.Get("email")
	user.ExpiresAt = token.Expiration.Unix()
	user.LogedAt = token.IssuedAt
	id, errID := strconv.ParseUint(token.Get("id"), 10, 32)

	if user.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("Expired token")
	}

	if errID != nil {
		return nil, fmt.Errorf("ERROR parsing the ID from the token. \nMore Information: %v", errID)
	}
	user.UserID = uint32(id)

	return user, nil
}

func (p *pasetoService) ParseToken(r *http.Request) string {
	return r.Header.Get("Authorization")
}

func (p *pasetoService) FormatSpecifics(s string) string {
	return s
}
