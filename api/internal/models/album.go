package models

import (
	"database/sql"
	"fmt"
	"koffee/internal/db"
	"koffee/pkg/formatter"
	"koffee/pkg/logger"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
)

type albumRepository struct {
	db        *sqlx.DB
	imageRepo RepositoryImages
}

// Album Database model
type Album struct {
	Name string `db:"name"`
	ID   uint32 `db:"id"`
	// NOTE(GABI): If we would need to implement or own see: https://gist.github.com/jmoiron/6979540
	Artists     pq.StringArray `db:"artists"`
	CoverURL    sql.NullString `db:"coverurl"`
	Description string         `db:"description"`
	Published   bool           `db:"published"`
	UploadDate  uint64         `db:"uploaddate"`
	PublishDate uint64         `db:"publishdate"`
	// Optional field when getting various albums
	// This is the total albums that the database found
	Fullcount uint64 `db:"fullcount"`
	// todo: (GABI) Discuss more fields
	// OPTIONAL FIELD
	UserID uint32 `db:"userid"`
	// OPTIONAL FIELD
	ArtistNames []string
	// OPTIONAL FIELD
	ArtistName  string `db:"artistname"`
	CoverImage  Image  `json:"cover_image"`
	HeaderImage Image  `json:"header_image"`
	// Invitation microservices
	EmailCreator string `json:"email_creator"`
}

// AlbumError model error
type AlbumError struct {
	UserID      string `json:user_id,omitempty`
	Name        string `json:"name,omitempty"`
	ID          string `json:"id,omitempty"`
	Artists     string `json:"artists,omitempty"`
	Description string `json:"description,omitempty"`
	Internal    string `json:"internal,omitempty"`
}

var albumSchema = `
	CREATE TABLE albums (
		name text,
		id SERIAL,
		artists text[],
		coverurl text NULL,
		description text,
		published BOOLEAN,
		uploaddate integer,
		publishdate integer
	)
`

// InitializeAlbums .
func InitializeAlbums(db *sqlx.DB, repoImages RepositoryImages) RepositoryAlbums {
	tx := db.MustBegin()
	_, _ = tx.Exec(albumSchema)
	_ = tx.Commit()
	return &albumRepository{db: db, imageRepo: repoImages}
}

func mapAlbumsToID(albums []Album) []uint32 {
	ints := make([]uint32, len(albums))
	for i := range albums {
		ints[i] = albums[i].ID
	}
	return ints
}

func (r *albumRepository) GetAlbumsImages(albums []Album) []Album {
	// GetProfilesImages updates the profiles array with the images, pass isSorted=false to sort byID
	images, err := r.imageRepo.GetImagesFromIDs([]ImageTypes{CoverImage, HeaderImage}, mapAlbumsToID(albums)...)
	if err != nil {
		return albums
	}
	fmt.Println(images)
	dict := Images(images).Zip()

	fmt.Println(dict)

	for idx := range albums {
		album := &albums[idx]

		imageDict, ok := dict[album.ID]

		if !ok {
			continue
		}

		if headerImage, ok := imageDict[HeaderImageAlbum]; ok {
			album.HeaderImage = headerImage
		}

		if albumImage, ok := imageDict[CoverImage]; ok {
			album.CoverImage = albumImage
		}
	}

	return albums

}

// CreateAlbum creates an album, but it's not published yet so the user can edit it (add songs or change cover or header)
func (r *albumRepository) CreateAlbum(userID uint32, name string, artistsIdentifiers []string, description string) (*Album, *AlbumError) {

	// Temporal:
	if len(artistsIdentifiers) != 1 {
		return nil, &AlbumError{Artists: "If you wanna add collaborators into the project you must invite them!"}
	}
	formatSearch, identifiers := formatter.Array(len(artistsIdentifiers), "username", artistsIdentifiers)
	logger.Log("album_model", formatSearch)
	tx := r.db.MustBegin()
	var profiles []Profile
	err := tx.Select(&profiles, fmt.Sprintf("SELECT * FROM profiles WHERE %s", formatSearch), identifiers...)
	logger.Log("album_model", fmt.Sprintf("SELECT * FROM profiles WHERE %s", formatSearch))
	if err != nil {
		return nil, &AlbumError{Internal: err.Error()}
	}
	if tx.Commit() != nil {
		return nil, &AlbumError{Internal: "error in tx"}
	}
	found := false
	for _, profile := range profiles {
		logger.Log("album_model", profile.UserID, userID)
		if profile.UserID == userID {
			found = true
		}
	}

	if !found {
		return nil, &AlbumError{UserID: "User making the request isn't in the artist array"}
	}

	if len(profiles) != len(artistsIdentifiers) {
		// return error
		return nil, &AlbumError{Artists: "Some artists don't exist!"} // ..
	}
	// !! Careful, maybe integer overflow
	lastID := -1
	tx = r.db.MustBegin()
	err = tx.QueryRowx("INSERT INTO albums (name, artists, coverurl, description, published, uploaddate) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", name, pq.StringArray(artistsIdentifiers), "", description, false, time.Now().Unix()).Scan(&lastID)
	if err != nil {
		// return error
		fmt.Println(err.Error())
		return nil, &AlbumError{Internal: "Error inserting"}
	}
	tx.Commit()
	if lastID == -1 {
		return nil, &AlbumError{ID: "Error inserting"}
	}
	return &Album{Name: name, Artists: artistsIdentifiers, Description: description, CoverURL: sql.NullString{String: ""}, Published: false, ID: uint32(lastID)}, nil
}

// GetAlbumByID returns an album by an ID (Not really safe)
func (r *albumRepository) GetAnyAlbumByID(id uint32) (*Album, *AlbumError) {
	var album []Album
	return r.getAlbum(album, "SELECT * FROM albums WHERE id=$1", id)
}

// GetPublicAlbumByID returns an album by its ID and makes sure its published. Also returns the artist full names
func (r *albumRepository) GetPublicAlbumByID(id uint32) (*Album, *AlbumError) {
	var album []Album
	albums, err := r.getAlbums(album, "SELECT albums.name, id, artists, coverurl, albums.description, published, uploaddate, profiles.name AS artistname FROM albums, profiles WHERE id=$1 AND published=TRUE AND profiles.username=ANY(albums.artists)", id)

	if err != nil {
		return nil, err
	}

	if len(albums) <= 0 {
		return nil, nil
	}

	artistNames := make([]string, len(albums))

	for i, album := range albums {
		artistNames[i] = album.Name
	}

	albums[0].ArtistNames = artistNames

	return &albums[0], err
}

// GetPrivateAlbumByID returns an album by ID and makes sure that it's owned by the passed userID and its private
func (r *albumRepository) GetPrivateAlbumByID(id uint32, userID uint32) (*Album, *AlbumError) {
	var album []Album
	return r.getAlbum(album, "SELECT albums.name AS name, id, artists, coverurl, albums.description, published, uploaddate FROM albums, profiles WHERE profiles.username=ANY(albums.artists) AND albums.id=$2 AND albums.published=FALSE AND profiles.userid=$1 GROUP BY albums.name, albums.id, albums.artists, albums.coverurl, albums.description, albums.published, albums.published, albums.uploaddate", userID, id)
}

// GetAlbumOwnedByID returns an album owned by the user of userID
func (r *albumRepository) GetAlbumOwnedByID(id uint32, userID uint32) (*Album, *AlbumError) {
	var album []Album
	return r.getAlbum(album, "SELECT albums.name AS name, id, artists, coverurl, albums.description, published, uploaddate FROM albums, profiles WHERE profiles.username=ANY(albums.artists) AND albums.id=$2 AND profiles.userid=$1 GROUP BY albums.name, albums.id, albums.artists, albums.coverurl, albums.description, albums.published, albums.published, albums.uploaddate", userID, id)
}

// GetAlbumOwnedByID returns an album owned by the user of userID
func (r *albumRepository) GetAlbumOwnedByIDPublish(id uint32, userID uint32, published bool) (*Album, *AlbumError) {
	var album []Album
	return r.getAlbum(album, "SELECT * FROM albums, profiles WHERE profiles.username=ANY(albums.artist) AND albums.id=$2 AND profiles.userid=$1 AND published=$3", userID, id, published)
}

// GetProfileAlbumsByUsername returns albums by username
func (r *albumRepository) GetProfileAlbumsByUsername(username string, published bool, afterID int, beforeID int, nItems int) ([]Album, *AlbumError) {
	var albums []Album
	orderBy := "albums.id > $3"
	id := afterID
	if beforeID != 0 && afterID == 0 {
		orderBy = "albums.id < $3"
		id = beforeID
	}
	s := fmt.Sprintf("SELECT albums.name, id, artists, coverurl, albums.description, published, uploaddate, profiles.userid, count(*) OVER() AS fullcount FROM albums, profiles WHERE $1=ANY(artists) AND published=$2 AND %s ORDER BY id ASC LIMIT $4", orderBy)
	fmt.Println(s)
	return r.getAlbums(albums, s, username, published, id, nItems)
}

// GetAlbumsByUserID return all albums which were done by the user
func (r *albumRepository) GetAlbumsByUserID(userID uint32, afterID int, beforeID int, nItems int) ([]Album, *AlbumError) {
	var albums []Album
	orderBy := "albums.id > $2"
	id := afterID
	if beforeID != 0 && afterID == 0 {
		orderBy = "albums.id < $2"
		id = beforeID
	}
	s := fmt.Sprintf("SELECT albums.name, id, artists, coverurl, albums.description, published, uploaddate, profiles.userid, count(*) OVER() AS fullcount FROM albums, profiles WHERE userid=$1 AND profiles.username=ANY(albums.artists) AND %s ORDER BY albums.uploaddate ASC LIMIT $3", orderBy)
	return r.getAlbums(albums, s, userID, id, nItems)
}

// GetAlbumsByUserID return all albums which were done by the user
func (r *albumRepository) GetAlbumsByUserIDPublish(userID uint32, published bool, afterID int, beforeID int, nItems int) ([]Album, *AlbumError) {
	var albums []Album
	orderBy := "albums.id > $3"
	id := afterID
	if beforeID != 0 && afterID == 0 {
		orderBy = "albums.id < $3"
		id = beforeID
	}
	s := fmt.Sprintf("SELECT albums.name, id, artists, coverurl, albums.description, published, uploaddate, profiles.userid, count(*) OVER() AS fullcount FROM albums, profiles WHERE userid=$1 AND profiles.username=ANY(albums.artists) AND published=$2 AND %s ORDER BY albums.uploaddate ASC LIMIT $4", orderBy)
	return r.getAlbums(albums, s, userID, published, id, nItems)
}

func (r *albumRepository) getAlbums(albums []Album, query string, params ...interface{}) ([]Album, *AlbumError) {
	tx := r.db.MustBegin()
	err := tx.Select(&albums, query, params...)
	if err != nil {
		return nil, &AlbumError{Internal: err.Error()}
	}
	errTx := tx.Commit()
	if errTx != nil {
		return nil, &AlbumError{Internal: errTx.Error()}
	}
	return albums, nil
}

func (r *albumRepository) getAlbum(albums []Album, query string, params ...interface{}) (*Album, *AlbumError) {
	tx := r.db.MustBegin()
	err := tx.Select(&albums, query, params...)
	if err != nil {
		return nil, &AlbumError{Internal: err.Error()}
	}
	if len(albums) <= 0 {
		return nil, nil
	}
	return &albums[0], nil
}

func (r *albumRepository) NewCollaborators(handlers []string, id uint32) *Album {
	tx := r.db.MustBegin()

	var album Album
	err := tx.QueryRowx("UPDATE albums SET artists=$1 WHERE albums.id=$2 RETURNING albums.artists, albums.name, albums.description, albums.published, albums.coverurl, albums.uploaddate", pq.StringArray(handlers), id).Scan(&album.Artists, &album.Name, &album.Description, &album.Published, &album.CoverURL, &album.UploadDate)
	tx.Commit()
	if err != nil {
		log.Println(err)
		return nil
	}
	return &album
}

func (r *albumRepository) UpdateAlbum(userID uint32, albumID uint32, publish string, description string, name string, coverURL string) (*Album, *AlbumError) {
	tx := r.db.MustBegin()
	album := Album{}
	values := make([]interface{}, 0)
	str := strings.Builder{}
	values = formatter.IfTrueAdd(&str, publish != "", "published", publish == "true", values)
	// todo reset table
	// values = formatter.IfTrueAdd(&str, publish != "", "publishdate", time.Now().Unix(), values)
	values = formatter.IfTrueAdd(&str, description != "", "description", description, values)
	values = formatter.IfTrueAdd(&str, name != "", "name", name, values)
	values = formatter.IfTrueAdd(&str, coverURL != "", "coverurl", coverURL, values)
	values = append(values, userID)
	values = append(values, albumID)
	s := fmt.Sprintf("UPDATE albums SET %s FROM profiles WHERE profiles.userid=$%d AND profiles.username=ANY(albums.artists) AND albums.id=$%d RETURNING profiles.userid AS userid, albums.artists, albums.name, albums.description, albums.published, albums.coverurl, albums.uploaddate", str.String(), len(values)-1, len(values))
	err := tx.QueryRowx(s, values...).Scan(&album.UserID, &album.Artists, &album.Name, &album.Description, &album.Published, &album.CoverURL, &album.UploadDate)

	if err != nil {
		return nil, &AlbumError{Internal: err.Error()}
	}
	err = tx.Commit()
	if err != nil {
		return nil, &AlbumError{Internal: err.Error()}
	}
	if album.UserID != userID {
		return nil, &AlbumError{UserID: "Error finding an album"}
	}
	return &album, nil
}

/**
{"name":"The great escape8","description":"Yes","published":false,"uploaddate":1587324025,"songs":null,"artists":[{"name":"Gabriel Villalonga","username":"gabivlj","userid":2,"artist":true,"imageurl":null,"headerimageurl":null,"description":null,"location":null}],"images":null}*/
// AlbumJSON the album for unmarshaling Postgresql query
type AlbumJSON struct {
	ID          uint64        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Published   bool          `json:"published"`
	UploadDate  uint64        `json:"uploaddate"`
	Songs       []Song        `json:"songs,omitempty"`
	Artists     []ProfileJSON `json:"artists"`
	Images      []Image       `json:"images,omitempty"`
}

// AlbumDBJSON is a row
type AlbumDBJSON struct {
	Album types.JSONText `db:"album"`
}

// GetAlbumFull retrieves an album in a JSON manner
// TODO: Test
func (r *albumRepository) GetAlbumFull(albumID uint32, published bool) (*AlbumJSON, error) {
	var albums AlbumDBJSON
	tx := r.db.MustBegin()
	err := tx.Get(&albums, db.GetAlbumFullInformation, albumID, published)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving album %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("Error retrieving album %v", err)
	}
	var album AlbumJSON
	err = albums.Album.Unmarshal(&album)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling album %v", err)
	}
	return &album, nil
}

// TODO: Test
func (r *albumRepository) GetAlbumsFull(published bool, ids ...uint32) ([]AlbumJSON, error) {
	var albums AlbumDBJSON
	tx := r.db.MustBegin()
	err := tx.Get(&albums, db.GetAlbumFullInformation, pq.Array(ids), published)
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("Error retrieving album %v", err)
	}
	var album []AlbumJSON
	err = albums.Album.Unmarshal(&album)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling album %v", err)
	}
	return album, nil
}
