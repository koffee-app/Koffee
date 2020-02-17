package models

import (
	"koffee/pkg/countries"

	"github.com/jmoiron/sqlx"
)

// Driver model
type Driver struct {
	Fullname string `db:"fullname"`
	UserID   uint32 `db:"id"`
	Budget   uint32 `db:"budget"`
	// ProfileImageURL string `db:"image_url"`
	Country string `db:"country"`
}

// DriverError represents an error doing an action with a Driver
type DriverError struct {
	Fullname string `json:"full_name"`
	UserID   string `json:"id"`
	Internal string `json:"internal"`
	Country  string `json:"country"`
}

// CreateDriver creates a driver
func CreateDriver(db *sqlx.DB, userID uint32, fullname string, country string) (*Driver, *DriverError) {
	u := UserByID(userID, db)
	if _, exists := countries.Country(country); !exists {
		return nil, &DriverError{Country: "Country does not exist"}
	}
	if u == nil {
		return nil, &DriverError{UserID: "Error, user does not exist with that ID"}
	}

	t := db.MustBegin()
	e := t.QueryRowx("INSERT INTO drivers (fullname, id, country) VALUES ($1, $2, $3)", fullname, userID, country)
	if e != nil {
		return nil, &DriverError{Internal: "Error inserting into database the driver"}
	}
	t.Commit()
	return &Driver{UserID: userID, Fullname: fullname, Country: country}, nil
}
