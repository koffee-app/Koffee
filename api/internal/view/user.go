package view

import (
	"net/http"
	"time"

	"koffee/internal/models"
)

// UserResponse Json
type UserResponse struct {
	Email      string    `json:"email"`
	Token      string    `json:"token"`
	ExpiresAt  int64     `json:"expires_at"`
	LogedAt    time.Time `json:"loged_at"`
	NewAccount bool      `json:"new_account"`
	ID         uint32    `json:"id"`
}

// User User response
func User(w http.ResponseWriter, user *models.User) {
	u := UserResponse{Email: user.Email, Token: user.Token, LogedAt: user.LogedAt, ExpiresAt: user.SessionExpiresAt, NewAccount: user.NewAccount, ID: user.UserID}
	PrepareJSON(w, http.StatusOK)
	ReturnJSON(w, FormatJSON(u, "Success", http.StatusOK))
}
