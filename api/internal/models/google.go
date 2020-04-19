package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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
func (r *repoUsers) GoogleAuthentication(token string) (*GoogleAuthResponse, *GoogleErrorResponse, error) {
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
func (r *repoUsers) addUserGoogleWithEmailCheck(email string) (*User, *UserError) {
	doesExist, _, _ := r.UserExists(email, nil)
	if doesExist {
		return nil, &UserError{Email: "User already exists"}
	}
	return r.addUserGoogle(email)
}

func (r *repoUsers) addUserGoogle(email string) (*User, *UserError) {
	tx := r.db.MustBegin()
	var lastID int
	_ = tx.QueryRowx("INSERT INTO users (email, password, isgoogleaccount) VALUES ($1, $2, $3) RETURNING id", email, "", true).Scan(lastID)
	tx.Commit()
	t, _ := r.tokenizer.GenerateToken(email, uint32(lastID), 30)
	refreshToken, _ := r.tokenizer.GenerateToken(email, uint32(lastID), 3000)
	return &User{Email: email, Token: r.tokenizer.FormatSpecifics(t), IsGoogleAccount: true, NewAccount: true, Password: sql.NullString{String: ""}, SessionExpiresAt: 0, RefreshToken: r.tokenizer.FormatSpecifics(refreshToken)}, nil
}

// LogUserGoogle Logs user via google access_token
func (r *repoUsers) LogUserGoogle(token string) (*User, *UserError) {
	succres, errres, err := r.GoogleAuthentication(token)
	if err != nil {
		return nil, &UserError{Internal: err.Error()}
	}
	if errres != nil {
		return nil, &UserError{Token: errres.ErrorToken}
	}

	exists, u, err := r.UserExists(succres.Email, nil)

	if err != nil {
		return nil, &UserError{Internal: err.Error()}
	}

	if !exists {
		return r.addUserGoogle(succres.Email)
	}

	if u == nil {
		// TODO: (GABI) Delete panic
		panic("Something bad happened inside UserExists()")
	}

	// todo ??? this
	// if !u.IsGoogleAccount {
	// 	return nil, &UserError{Email: "User is not a Google account!"}
	// }

	// TODO Please refactor this into a function hoooly
	u.Token, _ = r.tokenizer.GenerateToken(succres.Email, uint32(u.UserID), 30)
	u.Token = r.tokenizer.FormatSpecifics(u.Token)
	u.RefreshToken, _ = r.tokenizer.GenerateToken(succres.Email, uint32(u.UserID), 3000)
	u.RefreshToken = r.tokenizer.FormatSpecifics(u.RefreshToken)

	u.NewAccount = false
	u.Password = sql.NullString{String: ""}

	return u, nil
}
