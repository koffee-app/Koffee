package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const schema = `
	CREATE TABLE profiles (
		name				 text,
		username 		 text,
		userid 	 		 integer,
		artist	 		 boolean,
		imageurl 	   text NULL,
		headerimageurl text NULL,
		description	 text NULL,
		location		 text NULL
	)
`

type profileFind struct {
	fields     []interface{}
	fieldNames string
}

// Profile model
type Profile struct {
	Name           string         `db:"name"`
	Username       string         `db:"username"`
	UserID         uint32         `db:"userid"`
	Artist         bool           `db:"artist"`
	ImageURL       sql.NullString `db:"imageurl"`
	Location       sql.NullString `db:"location"`
	HeaderImageURL sql.NullString `db:"headerimageurl"`
	Description    sql.NullString `db:"description"`
}

// ProfileError is an error that will be returned when there is an error.
type ProfileError struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	UserID   string `json:"id"`
	Internal string `json:"internal"`
}

// InitializeProfile initializes tables if necessary of profiles
func InitializeProfile(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.Exec(schema)
	tx.Commit()
}

// CreateProfile creates a profile
// Returns a profile if succesful else if there is an error it will be stored inside ProfileError
func CreateProfile(db *sqlx.DB, username, name string, id uint32, artist bool) (*Profile, *ProfileError) {
	if err := checkFieldsCreate(username, name); err != nil {
		return nil, err
	}
	profile := wrapP(getProfile(db, &Profile{UserID: id, Username: username}))
	if profile != nil {
		if profile.UserID == id {
			return nil, &ProfileError{UserID: fmt.Sprintf("%d UserID already exists", id)}
		}
		if profile.Username == username {
			return nil, &ProfileError{Username: fmt.Sprintf("%s username already exists", username)}
		}
		return nil, &ProfileError{Internal: "We found profiles but we really don't know how..."}
	}
	tx := db.MustBegin()
	_, q := tx.Exec("INSERT INTO profiles (username, name, userid, artist) VALUES ($1, $2, $3, $4) RETURNING userid", username, name, id, artist)
	if q != nil {
		fmt.Println(q.Error())
		return nil, &ProfileError{Internal: "Error inserting into database, check logs"}
	}
	if e := tx.Commit(); e != nil {
		fmt.Println(e.Error() + " in createprofile final queryrowx")
		return nil, &ProfileError{Internal: "Error inserting into database, check logs"}
	}
	return &Profile{Username: username, UserID: id, Artist: artist, Name: name}, nil
}

// GetProfileByUsername returns a profile by username
func GetProfileByUsername(db *sqlx.DB, username string) *Profile {
	profile := wrapP(getProfile(db, &Profile{Username: username}))
	return profile
}

// GetProfileByUserID returns a profile by userID
func GetProfileByUserID(db *sqlx.DB, userID uint32) *Profile {
	profile := wrapP(getProfile(db, &Profile{UserID: userID}))
	return profile
}

func wrapP(profile *[]Profile) *Profile {
	if len(*profile) > 0 {
		return &(*profile)[0]
	}
	return nil
}

// Gets a profile by username or userID, will use OR if it username and id are both used.
func getProfile(db *sqlx.DB, profile *Profile) *[]Profile {
	var profiles []Profile
	tx := db.MustBegin()
	var varToUse profileFind
	builder := strings.Builder{}
	if profile.Username != "" {
		builder.WriteString("(username=$1")
		varToUse.fields = append(varToUse.fields, profile.Username)
	}
	if profile.UserID != 0 {
		if profile.Username != "" {
			builder.WriteString(" OR userid=$2)")
		} else {
			builder.WriteString("(userid=$1)")
		}
		varToUse.fields = append(varToUse.fields, profile.UserID)
	} else {
		builder.WriteByte(')')
	}
	varToUse.fieldNames = builder.String()
	query := fmt.Sprintf("SELECT * FROM profiles WHERE %s", varToUse.fieldNames)
	fmt.Println(query, varToUse.fields)
	e := tx.Select(&profiles, query, varToUse.fields...)
	if e != nil {
		fmt.Printf("Error getting a profile: " + e.Error())
		return &profiles
	}
	err := tx.Commit()
	if err != nil {
		fmt.Println("Error getting a profile in getProfile()")
		return &profiles
	}
	return &profiles
}

func checkFieldsCreate(username, name string) *ProfileError {
	if strings.Trim(username, " ") == "" || len(username) >= 20 {
		return &ProfileError{Username: "Username must be less than 21 characters or not empty"}
	}

	if strings.Trim(name, " ") == "" || len(name) >= 30 {
		return &ProfileError{Name: "The name must be less than 31 characters or not empty"}
	}
	return nil
}
