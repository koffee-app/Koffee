package middleware

import "net/http"

import "koffee/internal/models"

type userKey uint8

// UserContextKey Use this for getting from the context a user.
var UserContextKey userKey = 1

// GetUser Gets the user of the request.
func GetUser(r *http.Request) *models.User {
	u, _ := r.Context().Value(UserContextKey).(*models.User)
	return u
}
