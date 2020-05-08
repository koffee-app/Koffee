package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

const schemaSongs = `
	CREATE TABLE songs (
		name text,
		albumid integer,
		id SERIAL,
		duration integer,
		paths JSONB
	)
`

// SongModel represents a song in the database
type SongModel struct {
	ArtistFeatures types.JSONText `db:"artists"`
	Name           sql.NullString `db:"name"`
	ID             uint32         `db:"id"`
	AlbumID        uint32         `db:"albumid"`
	duration       uint64         `db:"duration"`
	paths          types.JSONText `db:"paths"`
}

type songRepo struct {
	db *sqlx.DB
}

// InitializeSongRepository Initializes the song repo
func InitializeSongRepository(db *sqlx.DB) songRepo {

	return songRepo{db: db}
}

// InsertSong inserts new song ( TEST )
func (s *songRepo) InsertSong(name string, albumID uint32) {

}
