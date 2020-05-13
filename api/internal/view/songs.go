package view

import (
	"koffee/internal/models"
	"net/http"
)

// SongsJSON contains the response data of the songs
type SongsJSON struct {
	Songs []models.Song `json:"songs"`
}

// Songs view
func Songs(w http.ResponseWriter, songs []models.Song) {
	s := SongsJSON{Songs: songs}
	PrepareJSON(w, http.StatusAccepted)
	data := JSON(s, "Success", http.StatusAccepted)
	ReturnJSON(w, data)
}

func Song(w http.ResponseWriter, song *models.Song) {
	PrepareJSON(w, http.StatusAccepted)
	data := JSON(song, "Success", http.StatusAccepted)
	ReturnJSON(w, data)
}
