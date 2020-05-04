package models

import (
	"database/sql"
	"fmt"
	"koffee/pkg/formatter"
	"koffee/pkg/logger"
	"log"
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

	ProfileImage Image `json:"profile_image"`
	HeaderImage  Image `json:"header_image"`
}

// Profiles arr
type Profiles []Profile

// Zip zips the profiles like {1: {...profile...}}
func (p Profiles) Zip() map[uint32]Profile {
	dict := make(map[uint32]Profile, p.Len())
	for _, profile := range p {
		dict[profile.UserID] = profile
	}
	return dict
}

// Len Sorting impl.
func (p Profiles) Len() int {
	return len(p)
}

// Swap is a order swap
func (p Profiles) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Less is a ordering func
func (p Profiles) Less(i, j int) bool {
	return p[i].UserID < p[j].UserID
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
	db        *sqlx.DB
	imageRepo RepositoryImages
}

// InitializeProfile initializes tables if necessary of profiles
func InitializeProfile(db *sqlx.DB, imageRepo RepositoryImages) RepositoryProfiles {
	tx := db.MustBegin()
	tx.Exec(schema)
	tx.Commit()
	return &repoProfiles{db: db, imageRepo: imageRepo}
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

func (r *repoProfiles) GetImage(profile *Profile) *Profile {
	imgs, _ := r.imageRepo.GetImagesSameID(profile.UserID, ProfileImage, HeaderImage)
	if len(imgs) > 1 {
		profile.ProfileImage = imgs[0]
		profile.HeaderImage = imgs[1]
	}
	if len(imgs) == 1 && imgs[0].Type == getImageType(ProfileImage) {
		profile.ProfileImage = imgs[0]
	} else if len(imgs) == 1 && imgs[0].Type == getImageType(HeaderImage) {
		profile.HeaderImage = imgs[0]
	}
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

func mapProfilesToID(profiles []Profile) []uint32 {
	ints := make([]uint32, len(profiles))
	for i := range profiles {
		ints[i] = profiles[i].UserID
	}
	return ints
}

func (r *repoProfiles) GetProfilesByIDs(profileIDs []uint32) []Profile {
	tx := r.db.MustBegin()
	// (userid=$1 OR userid=$2 ... OR userid=$N) - [1, 2, ..., n]
	query, arr := formatter.ArrayUint32(len(profileIDs), "userid", profileIDs)
	var profiles []Profile
	err := tx.Select(&profiles, fmt.Sprintf("SELECT * FROM profiles WHERE %s", query), arr...)
	if err != nil {
		log.Println(err)
		return []Profile{}
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return []Profile{}
	}

	// Zip profiles into a HashTable so we have O(1) accessing
	zippedProfiles := Profiles(profiles).Zip()

	// Order profiles as they were passed to the function
	for idx := range profileIDs {
		profiles[idx] = zippedProfiles[profileIDs[idx]]
		// if there are some profiles that weren't found in the Select query
		if idx+1 == len(profiles) {
			break
		}
	}

	return profiles
}

// GetProfilesImages updates the profiles array with the images, pass isSorted=false to sort byID
func (r *repoProfiles) GetProfilesImages(profiles []Profile) []Profile {
	images, err := r.imageRepo.GetImagesFromIDs([]ImageTypes{ProfileImage, HeaderImage}, mapProfilesToID(profiles)...)
	if err != nil {
		return profiles
	}

	dict := Images(images).Zip()

	for idx := range profiles {
		profile := &profiles[idx]

		imageDict, ok := dict[profile.UserID]

		if !ok {
			continue
		}

		if headerImage, ok := imageDict[HeaderImage]; ok {
			profile.HeaderImage = headerImage
		}

		if profileImage, ok := imageDict[ProfileImage]; ok {
			profile.ProfileImage = profileImage
		}
	}

	return profiles
}

// GetProfiles Gets profiles
func (r *repoProfiles) GetProfiles(profile *Profile, useArtistSearch bool, limit int) *[]Profile {
	var profiles []Profile
	tx := r.db.MustBegin()
	var varToUse profileFind
	builder := strings.Builder{}

	// (USERNAME=username
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

	// OR USERID=...
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

	// AND ARTIST=...
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
