package models

import (
	"database/sql"
	"fmt"
	"koffee/pkg/formatter"
	"koffee/pkg/logger"
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

// GetSingleProfile returns a single profile if it found it else nil
func (p *Profile) getSingleProfile(r *repoProfiles, useArtist bool) *Profile {
	return wrapP(r.GetProfiles(p, useArtist, 1))
}

type repoProfiles struct {
	db *sqlx.DB
}

// InitializeProfile initializes tables if necessary of profiles
func InitializeProfile(db *sqlx.DB) RepositoryProfiles {
	tx := db.MustBegin()
	tx.Exec(schema)
	tx.Commit()
	return &repoProfiles{db: db}
}

func (r *repoProfiles) SingleProfile(p *Profile, useArtist bool) *Profile {
	return p.getSingleProfile(r, useArtist)
}

// CreateProfile creates a profile
// Returns a profile if succesful else if there is an error it will be stored inside ProfileError
func (r *repoProfiles) CreateProfile(username, name string, id uint32, artist bool) (*Profile, *ProfileError) {
	if err := checkFieldsCreate(username, name); err != nil {
		return nil, err
	}
	profile := wrapP(r.GetProfiles(&Profile{UserID: id, Username: username}, false, 1))
	if profile != nil {
		if profile.UserID == id {
			return nil, &ProfileError{UserID: fmt.Sprintf("%d UserID already exists", id)}
		}
		if profile.Username == username {
			return nil, &ProfileError{Username: fmt.Sprintf("%s username already exists", username)}
		}
		return nil, &ProfileError{Internal: "We found profiles but we really don't know how..."}
	}
	tx := r.db.MustBegin()
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
func (r *repoProfiles) GetProfileByUsername(username string) *Profile {
	profile := wrapP(r.GetProfiles(&Profile{Username: username}, false, 1))
	return profile
}

// GetSingleProfile returns a single profile if it found it else nil
func (r *repoProfiles) GetSingleProfile(profile *Profile, useArtist bool) *Profile {
	return wrapP(r.GetProfiles(profile, useArtist, 1))
}

// GetProfileByUserID returns a profile by userID
func (r *repoProfiles) GetProfileByUserID(userID uint32) *Profile {
	profile := wrapP(r.GetProfiles(&Profile{UserID: userID}, false, 1))
	return profile
}

func wrapP(profile *[]Profile) *Profile {
	if len(*profile) > 0 {
		return &(*profile)[0]
	}
	return nil
}

// GetProfiles Gets profiles
// !FIXME this is totally overengineered
func (r *repoProfiles) GetProfiles(profile *Profile, useArtistSearch bool, limit int) *[]Profile {
	var profiles []Profile
	tx := r.db.MustBegin()
	var varToUse profileFind
	builder := strings.Builder{}

	formatter.FormatWhereQuery(
		profile.Username != "",
		len(varToUse.fields),
		"username",
		"",
		&builder,
		false,
		func() {
			varToUse.fields = append(varToUse.fields, profile.Username)
		},
	)

	formatter.FormatWhereQuery(
		profile.UserID != 0,
		len(varToUse.fields),
		"userid",
		"OR",
		&builder,
		false,
		func() {
			varToUse.fields = append(varToUse.fields, profile.UserID)
		},
	)

	formatter.FormatWhereQuery(
		useArtistSearch,
		len(varToUse.fields),
		"artist",
		"AND",
		&builder,
		true,
		func() {
			varToUse.fields = append(varToUse.fields, profile.Artist)
		},
	)

	varToUse.fieldNames = builder.String()
	query := fmt.Sprintf("SELECT * FROM profiles WHERE %s", varToUse.fieldNames)
	logger.Log("profile.go in CreateProfile()", query, varToUse.fields)
	e := tx.Select(&profiles, fmt.Sprintf("%s LIMIT %d", query, limit), varToUse.fields...)
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

// UpdateProfile .
func (r *repoProfiles) UpdateProfile(username string, description string, artist string, id uint32, name string) (*Profile, *ProfileError) {
	if username != "" && r.GetProfileByUsername(username) != nil {
		return nil, &ProfileError{Username: "Username already exists!"}
	}
	tx := r.db.MustBegin()
	profile := Profile{}
	values := make([]interface{}, 0)
	str := strings.Builder{}
	values = formatter.IfTrueAdd(&str, name != "", "name", name, values)
	values = formatter.IfTrueAdd(&str, username != "", "username", username, values)
	values = formatter.IfTrueAdd(&str, description != "", "description", description, values)
	values = formatter.IfTrueAdd(&str, artist != "", "artist", artist == "true", values)
	s := fmt.Sprintf("UPDATE profiles SET %s WHERE userid=$%d RETURNING profiles.userid, profiles.username, profiles.description, profiles.artist, profiles.age, profile.name", str.String(), len(values)-1)
	row := tx.QueryRowx(s, values...).Scan(&profile.UserID, &profile.Username, &profile.Description, &profile.Artist)
	if row != nil && row.Error() != "" {
		return nil, &ProfileError{Internal: row.Error()}
	}
	if err := tx.Commit(); err != nil {
		return nil, &ProfileError{Internal: err.Error()}
	}

	return &profile, nil
}
