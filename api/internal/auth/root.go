package auth

import (
	"net/http"
	"time"
)

// IUser is users interface for tokens
type IUser interface {
	Information() (string, time.Time, uint32, int64)
}

// Token interface that needs to be implemented, only for generating and verifying tokens
type Token interface {
	// Generates a token (email, ID, duration in Minutes) (token, error)
	GenerateToken(string, uint32, uint64) (string, error)
	// Checks if the token included in the request is valid (Only checks if its valid by itself, doesnt make a Database read)
	TokenValid(*http.Request) (IUser, error)
	// Verifies token (Only checks if its valid by itself, doesnt make a Database read)
	VerifyToken(tokenStr string) (IUser, error)
	// Gets token from request
	ParseToken(r *http.Request) string

	FormatSpecifics(string) string
}
