package middleware

import (
	"context"

	"github.com/julienschmidt/httprouter"

	"net/http"

	"koffee/internal/auth"
	"koffee/internal/models"
	"koffee/internal/view"
)

// JwtAuthentication Middleware for protected routes
func JwtAuthentication(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		user, err := auth.TokenValid(r)
		if err != nil {
			view.ErrorAuthentication(w, map[string]string{"why": err.Error()})
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, models.UserJWTToUser(user))
		r = r.WithContext(ctx)
		next(w, r, p)
	}
}
