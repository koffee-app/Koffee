package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const schema = `
	CREATE TABLE profiles (
		username 		 text,
		userID 	 		 integer,
		artist	 		 boolean,
		imageurl 	 text NULL,
		description	 text NULL,
		age					 integer NULL
	)
`

type profileFind struct {
	fields     []interface{}
	fieldNames string
}

// Profile model
type Profile struct {
	Username    string         `db:"username"`
	UserID      uint32         `db:"id"`
	Artist      bool           `db:"artist"`
	Age         sql.NullInt64  `db:"age"`
	ImageURL    sql.NullString `db:"imageurl"`
	Description sql.NullString `db:"description"`
}

// ProfileError is an error that will be returned when there is an error.
type ProfileError struct {
	Username string `json:"username"`
	UserID   string `json:"id"`
	Internal string `json:"internal"`
}

// CreateProfile creates a profile
// Returns a profile if succesful else if there is an error it will be stored inside ProfileError
func CreateProfile(db *sqlx.DB, username string, id uint32, artist bool) (*Profile, *ProfileError) {
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
	q := tx.QueryRowx("INSERT INTO profiles (username, id, artist) VALUES ($1, $2, $3)", username, id, artist)
	if q.Err() != nil {
		fmt.Println(q.Err().Error())
		return nil, &ProfileError{Internal: "Error inserting into database, check logs"}
	}
	if e := tx.Commit(); e != nil {
		fmt.Println(e.Error())
		return nil, &ProfileError{Internal: "Error inserting into database, check logs"}
	}
	return &Profile{Username: username, UserID: id, Artist: artist}, nil
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
		builder.WriteString("($1=$2")
		varToUse.fields = append(varToUse.fields, "username", profile.Username)
	}
	if profile.UserID != 0 {
		if profile.Username != "" {
			builder.WriteString("OR $3=$4)")
		} else {
			builder.WriteString("$1=$2)")
		}
		varToUse.fields = append(varToUse.fields, "id", profile.UserID)
	} else {
		builder.WriteByte(')')
	}
	varToUse.fieldNames = builder.String()
	e := tx.Select(&profiles, "SELECT * FROM profiles WHERE "+varToUse.fieldNames, varToUse.fields...)
	if e != nil {
		fmt.Printf("Error getting a profile: " + e.Error())
		return &profiles
	}
	tx.Commit()
	return &profiles
}
