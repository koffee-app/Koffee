package models

import (
	"database/sql"
	"koffee/pkg/formatter"
	"koffee/pkg/logger"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Album Database model
type Album struct {
	Name string `db:"name"`
	ID   uint32 `db:"id"`
	// NOTE(GABI): If we would need to implement or own see: https://gist.github.com/jmoiron/6979540
	Artists     pq.StringArray `db:"artists"`
	CoverURL    sql.NullString `db:"coverurl"`
	Description string         `db:"description"`
	// todo: (GABI) Discuss more fields
}

var albumSchema = `
	CREATE TABLE albums (
		name text,
		id SERIAL,
		artists text[],
		coverurl text NULL,
		description text
	)
`

// CreateAlbum creates an album (todo docs)
func CreateAlbum(db *sqlx.DB, name string, artistsIdentifiers []string, description string) {
	logger.Log("album_model", "NOT IMPLEMENTED")
	panic("not implemented")
	formatSearch := formatter.Array(len(artistsIdentifiers), "name")
	logger.Log("album_model", formatSearch)
	tx := db.MustBegin()
	var profiles Profile[]
	tx.Select(&profiles, fmt.Sprintf("SELECT FROM profile WHERE %s", formatSearch), ...artistsIdentifiers)
	if len(*profiles) != len(artistIdentifiers) {
		// return error
		return 
	}
	
}
