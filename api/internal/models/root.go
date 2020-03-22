package models

import "github.com/jmoiron/sqlx"

// Initialize every table
func Initialize(db *sqlx.DB) {
	InitializeDrivers(db)
	InitializeUsers(db)
	InitializeProfile(db)
}
