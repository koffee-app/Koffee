package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"koffee/internal/auth"
	"net/http"

	"github.com/jmoiron/sqlx"
)

var url = `https://oauth2.googleapis.com/tokeninfo?access_token=token`

// GoogleAuthResponse 200 OK Google response
// JSON Response example:
// "iss": "https://accounts.google.com",
// "sub": "110169484474386276334",
// "azp": "1008719970978-hb24n2dstb40o45d4feuo2ukqmcc6381.apps.googleusercontent.com",
// "aud": "1008719970978-hb24n2dstb40o45d4feuo2ukqmcc6381.apps.googleusercontent.com",
// "iat": "1433978353",
// "exp": "1433981953",
// ****** These seven fields are only included when the user has granted the "profile" and
// "email" OAuth scopes to the application.
// "email": "testuser@gmail.com",
// "email_verified": "true",
// "name" : "Test User",
type GoogleAuthResponse struct {
	ISS   string `json:"iss"`
	Sub   string `json:"sub"`
	AZP   string `json:"azp"`
	AUD   string `json:"aud"`
	IAT   string `json:"iat"`
	EXP   string `json:"exp"`
	Email string `json:"email"`
}

// GoogleErrorResponse Error JSON for the google response
// Example: {
//   "error": "invalid_token",
//   "error_description": "Invalid Value"
// }
type GoogleErrorResponse struct {
	ErrorToken       string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// GoogleAuthentication which returns a GoogleAuthResponse or GoogleErrorResponse if the token is wrong,
// don't use this in highly frequented requests or middleware because it's gonna be expensive! Only for one time verif.
func GoogleAuthentication(token string) (*GoogleAuthResponse, *GoogleErrorResponse, error) {
	client := http.Client{}
	resp, err := client.Get(fmt.Sprintf(`https://oauth2.googleapis.com/tokeninfo?access_token=%s`, token))
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		var e GoogleErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return nil, nil, err
		}
		return nil, &e, nil
	}
	var rsucc GoogleAuthResponse
	err = json.NewDecoder(resp.Body).Decode(&rsucc)
	if err != nil {
		return nil, nil, err
	}
	return &rsucc, nil, nil
}

// Adds a Google User, check if it exists
// todo Check more errors.
func addUserGoogle(db *sqlx.DB, email string) (*User, *UserError) {
	doesExist, _ := UserExists(email, nil, db)
	if doesExist {
		return nil, &UserError{Email: "User already exists"}
	}
	tx := db.MustBegin()
	var lastID int
	_ = tx.QueryRowx("INSERT INTO users (email, password, isgoogleaccount) VALUES ($1, none, $2) RETURNING id", email, true).Scan(lastID)
	tx.Commit()
	t, _ := auth.GenerateTokenJWT(email, uint32(lastID))
	expires := auth.ExpiresAt()
	return &User{Email: email, Token: fmt.Sprintf("Bearer %s", t), IsGoogleAccount: true, NewAccount: true, Password: sql.NullString{String: ""}, SessionExpiresAt: expires}, nil
}
