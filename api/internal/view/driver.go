package view

//! !!!!!! DEPRECATED

import (
	"koffee/internal/models"
	"net/http"
)

// DriverJSON /
type DriverJSON struct {
	UserID   uint32 `json:"id,omitempty"`
	Country  string `json:"country,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	Budget   uint32 `json:"budget,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

// Driver renders driver json view
func Driver(w http.ResponseWriter, d *models.Driver) {
	driver := DriverJSON{UserID: d.UserID, Country: d.Country, Fullname: d.Fullname, ImageURL: /* todo */ "", Budget: d.Budget}
	ReturnJSON(w, FormatJSON(driver, "Success", http.StatusOK))
}
