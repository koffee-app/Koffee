package models

import (
	"database/sql"
	"fmt"
	"koffee/pkg/countries"

	"github.com/jmoiron/sqlx"
)

//! !!!!!! DEPRECATED

const driverSchema = `
	CREATE TABLE drivers (
		fullname text,
		id			 integer,
		budget	 real,
		country  text,
		imageurl text NULL
	)
`

// Driver model
type Driver struct {
	Fullname        string         `db:"fullname"`
	UserID          uint32         `db:"id"`
	Budget          uint32         `db:"budget"`
	ProfileImageURL sql.NullString `db:"imageurl"`
	Country         string         `db:"country"`
}

// DriverError represents an error doing an action with a Driver
type DriverError struct {
	Fullname string `json:"full_name"`
	UserID   string `json:"id"`
	Internal string `json:"internal"`
	Country  string `json:"country"`
}

// InitializeDrivers Initializes table of drivers
func InitializeDrivers(db *sqlx.DB) {
	t := db.MustBegin()
	t.Exec(driverSchema)
	t.Commit()
}

// CreateDriver creates a driver
func CreateDriver(db *sqlx.DB, userID uint32, fullname, country string) (*Driver, *DriverError) {
	if _, exists := countries.Country(country); !exists {
		return nil, &DriverError{Country: fmt.Sprintf("Country code %s does not exist.", country)}
	}
	u := UserByID(userID, db)
	if u == nil {
		return nil, &DriverError{UserID: "Error, user does not exist with that ID"}
	}
	t := db.MustBegin()
	e := t.QueryRowx("INSERT INTO drivers (fullname, id, country, budget) VALUES ($1, $2, $3, $4)", fullname, userID, country, 0.0)
	if e.Err() != nil {
		fmt.Println(e.Err())
		return nil, &DriverError{Internal: "Error inserting into database the driver"}
	}
	t.Commit()
	return &Driver{UserID: userID, Fullname: fullname, Country: country}, nil
}

// GetDriverByID return a driver reference, if it doesn't exist returns nil.
// (GABI): Consideration, maybe add a explanation which tells the caller what
//				 happened
func GetDriverByID(db *sqlx.DB, userID uint32) *Driver {
	tx := db.MustBegin()
	drivers := []Driver{}
	err := tx.Select(&drivers, "SELECT * FROM drivers WHERE id=$1", userID)
	if err != nil {
		//  todo (GABI): handle error see Consideration note
		fmt.Println("ERROR in GetDriverByID:\n" + err.Error())
		return nil
	}
	if len(drivers) == 0 {
		return nil
	}
	return &(drivers[0])
}
