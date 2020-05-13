package models

import (
	"database/sql"
	"errors"
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
	Song           types.JSONText `db:"song"`
}

// Song model
type Song struct {
	ArtistFeatures []ProfileJSON `json:"artists"`
	Name           string        `json:"name"`
	ID             uint32        `json:"id"`
	AlbumID        uint32        `json:"albumid"`
	Duration       uint64        `json:"duration"`
}

// Songs is a songs array json
type Songs struct {
	Song []Song `json:"songs"`
}

type songRepo struct {
	db *sqlx.DB
}

// ErrSongNotFound is an error when the requested song/s are not in the database
var ErrSongNotFound = errors.New("Song or songs requested are not found in the database")

// InitializeSongRepository Initializes the song repo
func InitializeSongRepository(db *sqlx.DB) SongsRepository {
	return &songRepo{db: db}
}

// InsertSong inserts new song ( TEST )
func (s *songRepo) InsertSong(name string, albumID uint32) {

}

func (s *songRepo) GetSongsByID(albumID uint32) ([]Song, error) {
	tx := s.db.MustBegin()
	var song SongModel
	err := tx.Get(&song, db.GetSongsByAlbumIDQuery, albumID, false)

	if err == sql.ErrNoRows {
		return nil, ErrSongNotFound
	}

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

func (s *songRepo) GetSongByID(songID uint32) (*Song, error) {
	tx := s.db.MustBegin()
	var song SongModel
	err := tx.Get(&song, db.GetSongsByAlbumIDQuery, songID)

	if err == sql.ErrNoRows {
		return nil, ErrSongNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("Error with the transaction %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("Error with the transaction %v", err)
	}
	var songJSON Song
	if err := song.Song.Unmarshal(&songJSON); err != nil {
		return nil, fmt.Errorf("Error with the transaction %v", err)
	}
	return &songJSON, nil
}
