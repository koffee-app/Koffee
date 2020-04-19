package view

import (
	"koffee/internal/models"
	"net/http"
)

type AlbumJSON struct {
	Name        string        `json:"name,omitempty"`
	ID          uint32        `json:"id,omitempty"`
	Artists     []string      `json:"artists,omitempty"`
	CoverURL    string        `json:"cover_url,omitempty"`
	Description string        `json:"description,omitempty"`
	Published   bool          `json:"published"`
	UploadDate  uint64        `json:"upload_date,omitempty"`
	PublishDate uint64        `json:"publish_date,omitempty"`
	Fullcount   uint64        `json:"full_count,omitempty"`
	UserID      uint32        `json:"user_id,omitempty"`
	ArtistNames []string      `json:"artist_names,omitempty"`
	ArtistName  string        `json:"artist_name,omitempty"`
	Songs       []interface{} `json:"songs,omitempty"` // todo
}

// Album writes to the client the album
func Album(w http.ResponseWriter, a *models.Album) {
	album := AlbumJSONFromModel(a)
	PrepareJSON(w, http.StatusOK)
	ReturnJSON(w, JSON(album, "Success", http.StatusOK))
}

// Albums writes to the client the albums
func Albums(w http.ResponseWriter, a []models.Album) {
	albums := make([]*AlbumJSON, len(a))
	for idx := range a {
		albums[idx] = AlbumJSONFromModel(&a[idx])
	}
	PrepareJSON(w, http.StatusOK)
	ReturnJSON(w, JSON(albums, "Success", http.StatusOK))
}

// AlbumError renders the album error
func AlbumError(w http.ResponseWriter, a *models.AlbumError) {
	if a.Internal != "" {
		Error(w, "Internal error, more details in data!", http.StatusNotFound, a)
		return
	}
	PrepareJSON(w, http.StatusBadRequest)
	ReturnJSON(w, JSON(a, "Error, see data for more details", http.StatusBadRequest))
}

// AlbumJSONFromModel returns the AlbumJSON corresponding to model album
func AlbumJSONFromModel(a *models.Album) *AlbumJSON {
	return &AlbumJSON{Name: a.Name, ID: a.ID, Artists: a.Artists, Description: a.Description, Published: a.Published, UploadDate: a.UploadDate, PublishDate: a.PublishDate, Fullcount: a.Fullcount, UserID: a.UserID, ArtistNames: a.ArtistNames, ArtistName: a.ArtistName}
}
