// todo: VerifyUser(email string, encryptedPassword string) -> if it exists and its a google account return bad login
// todo: VerifyUserGoogle(email string, token string) -> for checking the user in the db, if it exists and its not a google account change it as one or inform the user. If it
//																												does not exist, create it as a google acc and return as a NewAccount
// todo: GetUserByEmail(email string, getEncryptedPassword bool)
// todo: GetUserByID(id uint32, getEncryptedPassword bool)
// todo: AddUser(email string, password string, isGoogleAccount bool) -> if it's google account check that password is empty

package models

import (
	"database/sql"
	"fmt"
	"koffee/internal/auth"
	"time"

	// "koffee/internal/auth"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

const userSchema = `
	CREATE TABLE users (
		email text,
		id		SERIAL,
		password text NULL,
		isGoogleAccount bool NULL
	);
`

// User user model. We will connect the UserID with a profile model that may have more information of the user
type User struct {
	Email            string         `db:"email"`
	UserID           uint32         `db:"id"`
	Password         sql.NullString `db:"password"`
	IsGoogleAccount  bool           `db:"isgoogleaccount"`
	NewAccount       bool
	Token            string
	LogedAt          time.Time
	SessionExpiresAt int64
}

// UserError error body
type UserError struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Internal string `json:"internal"`
}

// Initialize initializes table
func Initialize(db *sqlx.DB) {
	t := db.MustBegin()
	t.Exec(userSchema)
	t.Commit()
}

// Message informing what happened.
func (u *UserError) Message() string {
	return `Error with authenticating user, see more information in data fields.`
}

// AddUser tries to insert a new User into the database via
func AddUser(db *sqlx.DB, email, password string, isGoogleAccount bool) (*User, *UserError) {
	if !isGoogleAccount {
		return addUserNoGoogle(db, email, password)
	}
	return nil, nil
}

// LogUserNotGoogle Tries to log in via email password way, if the stored user is a google acc this will
// return an error.
func LogUserNotGoogle(db *sqlx.DB, email, password string) (*User, *UserError) {
	u := User{}
	users := &[]User{}
	t := db.MustBegin()
	email = strings.ToLower(email)
	e := t.Select(users, "SELECT * FROM users WHERE email=$1 LIMIT 1", email)

	if e != nil {
		return nil, &UserError{Email: "Incorrect credentials.", Internal: e.Error()}
	}

	if len(*users) == 0 {
		return nil, &UserError{Email: "Incorrect credentials"}
	}

	u = (*users)[0]

	if u.Email == "" || u.IsGoogleAccount {
		return nil, &UserError{Email: "Incorrect credentials"}
	}

	pass := u.Password.String
	if err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password)); err != nil {
		return nil, &UserError{Password: "Incorrect credentials."}
	}
	u.Token, e = auth.GenerateTokenJWT(email, u.UserID)
	if e != nil {
		return nil, &UserError{Email: "Incorrect user credentials?", Internal: e.Error(), Password: "Error generating token"}
	}
	return &u, nil
}

// adds a user to the database taking into mind that it's not a google user
func addUserNoGoogle(db *sqlx.DB, email, password string) (*User, *UserError) {
	t := db.MustBegin()
	// sanitize email input.
	email = strings.ToLower(email)
	if userError := checkValuesForAdding(db.MustBegin(), email, password); userError != nil {
		return nil, userError
	}
	// start knowing that the inputs is sanitized and safe.
	p := encryptPassword(password)
	// lastID declaration for getting it from Scan
	var lastID int
	// insert the user and tget the lastID
	err := t.QueryRowx("INSERT INTO users (email, password, isgoogleaccount) VALUES ($1, $2, $3) RETURNING id", email, p, false).Scan(&lastID)
	if err != nil {
		return nil, &UserError{Internal: err.Error() + " line 72"}
	}
	// commit changes
	t.Commit()
	// generate token
	token, _ := auth.GenerateTokenJWT(email, uint32(lastID))
	// return newly generated user
	return &User{Password: sql.NullString{String: p}, IsGoogleAccount: false, Email: email, NewAccount: true, Token: fmt.Sprintf("Bearer %s", token), UserID: uint32(lastID)}, nil
}

// compares hashed and nonhashed passwords and if they are equal returns true
func areEqual(nonHashed, hashed string) bool {
	success := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(nonHashed))
	if success != nil {
		return false
	}
	return true
}

// check if the values for adding the user into the db are legit.
func checkValuesForAdding(tx *sqlx.Tx, email, password string) *UserError {
	e := &UserError{Internal: ""}
	bad := false

	if len(password) < 6 {
		e.Password = "Password length should be 6 or more"
		bad = true
	} else if len(password) > 30 {
		e.Password = "Password length should be 6 or more"
		bad = true
	}

	if len(email) > 320 {
		bad = true
		e.Email = "Invalid email address"
	} else if strings.Index(email, "@") < 0 {
		bad = true
		e.Email = "Email should include @"
	}

	if bad {
		return e
	}
	// Start checking if the email exists.
	u := &[]User{}
	// Check if email already exists.
	err := tx.Select(u, "SELECT email, password FROM users WHERE email=$1 LIMIT 1", email)
	if err != nil {
		fmt.Println(err.Error())
		e.Internal = err.Error()
		return e
	}
	tx.Commit()

	// If it does not exist all is fine.
	if len(*u) == 0 {
		return nil
	}

	e.Email = `User already exists!`
	return e
}

// Encrypts password
func encryptPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("Error hashing password.")
	}
	return string(hashed)
}

// UserJWTToUser .
func UserJWTToUser(u *auth.UserJWT) *User {
	return &User{Email: u.Email, UserID: u.UserID, LogedAt: u.LogedAt, SessionExpiresAt: u.ExpiresAt, NewAccount: false}
}

// todo: DeleteUser(email string)
// todo: DeleteUserViaID(userID uint32)
