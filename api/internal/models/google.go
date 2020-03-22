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

// GoogleAuthentication which returns a GoogleAuthResponse or GoogleErrorResponse if the token is wrong, sends a request to google servers to check if the token is fine
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

// Adds a Google User, checks if it exists in our db and if it does return error
// todo Check more errors.
func addUserGoogle(db *sqlx.DB, email string) (*User, *UserError) {
	doesExist, _, _ := UserExists(email, nil, db)
	if doesExist {
		return nil, &UserError{Email: "User already exists"}
	}
	tx := db.MustBegin()
	var lastID int
	_ = tx.QueryRowx("INSERT INTO users (email, password, isgoogleaccount) VALUES ($1, $2, $3) RETURNING id", email, "", true).Scan(lastID)
	tx.Commit()
	t, _ := auth.GenerateTokenJWT(email, uint32(lastID))
	expires := auth.ExpiresAt()
	return &User{Email: email, Token: fmt.Sprintf("Bearer %s", t), IsGoogleAccount: true, NewAccount: true, Password: sql.NullString{String: ""}, SessionExpiresAt: expires}, nil
}

// LogUserGoogle Logs user via google access_token
func LogUserGoogle(db *sqlx.DB, token string) (*User, *UserError) {
	succres, errres, err := GoogleAuthentication(token)
	if err != nil {
		return nil, &UserError{Internal: err.Error()}
	}
	if errres != nil {
		return nil, &UserError{Token: errres.ErrorToken}
	}

	exists, u, err := UserExists(succres.Email, nil, db)

	if err != nil {
		return nil, &UserError{Internal: err.Error()}
	}

	if !exists {
		return nil, &UserError{Email: "User does not exist!"}
	}

	if u == nil {
		// TODO: (GABI) Delete panic
		panic("Something bad happened inside UserExists()")
	}

	if !u.IsGoogleAccount {
		return nil, &UserError{Email: "User is not a Google account!"}
	}

	u.Token, _ = auth.GenerateTokenJWT(succres.Email, uint32(u.UserID))
	u.Token = auth.Bearify(u.Token)
	u.SessionExpiresAt = auth.ExpiresAt()
	u.NewAccount = false
	u.Password = sql.NullString{String: ""}

	return u, nil
}
