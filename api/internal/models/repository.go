package models

import (
	"github.com/jmoiron/sqlx"
)

// RepositoryUsers interface
type RepositoryUsers interface {
	AddUser(email, password string, isGoogleAccount bool) (*User, *UserError)
	LogUserNotGoogle(email, password string) (*User, *UserError)
	LogUserGoogle(token string) (*User, *UserError)
	GoogleAuthentication(token string) (*GoogleAuthResponse, *GoogleErrorResponse, error)
	UserExists(email string, tx *sqlx.Tx) (bool, *User, error)
	UserByID(id uint32) *User
	Refresh(refreshToken string) (*User, *UserError)
}

// RepositoryProfiles interface
type RepositoryProfiles interface {
	CreateProfile(username, name string, id uint32, artist bool) (*Profile, *ProfileError)
	SingleProfile(p *Profile, useArtist bool) *Profile
	GetProfileByUsername(username string) *Profile
	GetProfileByUserID(userID uint32) *Profile
	GetProfiles(profile *Profile, useArtistSearch bool, limit int) *[]Profile
}

// RepositoryAlbums interface
type RepositoryAlbums interface {
	CreateAlbum(userID uint32, name string, artistsIdentifiers []string, description string) (*Album, *AlbumError)
	// Returns an album which is published in the platform
	GetPublicAlbumByID(id uint32) (*Album, *AlbumError)
	GetPrivateAlbumByID(id uint32, userID uint32) (*Album, *AlbumError)
	GetAlbumOwnedByID(id uint32, userID uint32) (*Album, *AlbumError)
	// Doesn't care if it's public or not. You can make sure that the album is owned by the user checking the artists array NOTE: (We are gonna delete this i think)
	GetAnyAlbumByID(id uint32) (*Album, *AlbumError)
	// GetProfileAlbumsByUsername returns albums by username
	GetProfileAlbumsByUsername(username string, published bool, afterID int, beforeID int, nItems int) ([]Album, *AlbumError)
	GetAlbumsByUserID(userID uint32, afterID int, beforeID int, nItems int) ([]Album, *AlbumError)
	GetAlbumOwnedByIDPublish(id uint32, userID uint32, published bool) (*Album, *AlbumError)
	UpdateAlbum(userID uint32, albumID uint32, publish string, description string, name string, coverURL string) (*Album, *AlbumError)
	GetAlbumsByUserIDPublish(userID uint32, published bool, afterID int, beforeID int, nItems int) ([]Album, *AlbumError)
}
