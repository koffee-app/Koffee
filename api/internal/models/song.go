package models

import (
	"database/sql"
	"fmt"
	"koffee/internal/db"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

const schemaSongs = `
	CREATE TABLE songs (
		name text,
		albumid integer,
		id SERIAL,
		duration integer,
		paths JSONB,
		artists text[]
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
	Songs          types.JSONText `db:"songs"`
}

// Song model
type Song struct {
	ArtistFeatures []ProfileJSON `json:"artists"`
	Name           string        `json:"name"`
	ID             uint32        `json:"id"`
	AlbumID        uint32        `json:"albumid"`
	Duration       uint64        `json:"duration"`
}

type Songs struct {
	Song []Song `json:"songs"`
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

func (s *songRepo) GetSongsByID(albumID uint32) ([]Song, error) {
	tx := s.db.MustBegin()
	var song SongModel
	err := tx.Get(&song, db.GetSongsByAlbumIDQuery, albumID, false)
	if err != nil {
		return nil, fmt.Errorf("Error with the transaction %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("Error with the transaction %v", err)
	}
	var songs Songs
	err = song.Songs.Unmarshal(&songs)
	if err != nil {
		return nil, fmt.Errorf("Error with the transaction %v", err)
	}
	return songs.Song, nil
}
