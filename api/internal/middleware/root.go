package middleware

import (
	"koffee/internal/auth"
	"koffee/internal/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type userKey uint8

// UserContextKey Use this for getting from the context a user.
var UserContextKey userKey = 1

var tokenService auth.Token

// Initialize initializes middleware
func Initialize(tokenizer auth.Token) {
	tokenService = tokenizer
}

// GetUser Gets the user of the request.
func GetUser(r *http.Request) *models.User {
	u, _ := r.Context().Value(UserContextKey).(*models.User)
	return u
}

// Authorization middleware
func Authorization(h httprouter.Handle) httprouter.Handle {
	return JwtAuthentication(h, tokenService)
}
