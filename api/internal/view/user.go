package view

import (
	"net/http"
	"time"

	"koffee/internal/models"
)

// UserResponse Json
type UserResponse struct {
	Email        string    `json:"email,omitempty"`
	Token        string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    int64     `json:"expires_at,omitempty"`
	LogedAt      time.Time `json:"loged_at,omitempty"`
	NewAccount   bool      `json:"new_account,omitempty"`
	ID           uint32    `json:"id,omitempty"`
}

// User User response
func User(w http.ResponseWriter, user *models.User) {
	if user.LogedAt.Year() == 1 {
		user.LogedAt = time.Now()
	}
	u := UserResponse{Email: user.Email, Token: user.Token, LogedAt: user.LogedAt, ExpiresAt: user.SessionExpiresAt, NewAccount: user.NewAccount, ID: user.UserID, RefreshToken: user.RefreshToken}
	PrepareJSON(w, http.StatusOK)
	ReturnJSON(w, JSON(u, "Success", http.StatusOK))
}
