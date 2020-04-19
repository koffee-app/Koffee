package models

import (
	"koffee/internal/auth"

	"github.com/jmoiron/sqlx"
)

/// Models will use a PostgreSQL database for the queries.
/// We won't do indexed search and cache here! That will be inside another folder

// Initialize every table
func Initialize(db *sqlx.DB, tokenService auth.Token) (RepositoryUsers, RepositoryProfiles, RepositoryAlbums) {
	// todo Reutrn all repos
	return InitializeUsers(db, tokenService), InitializeProfile(db), InitializeAlbums(db)
}
