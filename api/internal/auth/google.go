package auth

import (
	"net/http"
)

// UserGoogle Google user.
type UserGoogle struct {
	email string
	exp   uint32
}

// VerifyGoogleToken .
func VerifyGoogleToken(token string) error {
	return nil
}

// Google Gets the user google from the request (access_token)
func Google(r *http.Request) *UserGoogle {
	return nil
}

// GoogleToUser Gets the User from google model
// func GoogleToUser(u *UserGoogle) *models.User {
// 	// TODO: Get from DB the user
// 	return &models.User{UserID: 0, Email: u.email}
// }
