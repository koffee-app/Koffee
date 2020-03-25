package view

import (
	"koffee/internal/models"
	"net/http"
)

type profileJSON struct {
	Username       string `json:"username,omitempty"`
	Name           string `json:"name,omitempty"`
	UserID         uint32 `json:"id,omitempty"`
	Artist         bool   `json:"artist,omitempty"`
	Location       string `json:"location,omitempty"`
	HeaderImageURL string `json:"headerimageurl,omitempty"`
	ImageURL       string `json:"imageurl,omitempty"`
	Description    string `json:"description,omitempty"`
}

// Profile Renders profile JSON
func Profile(w http.ResponseWriter, profile *models.Profile) {
	PrepareJSON(w, http.StatusOK)
	ReturnJSON(w, FormatJSON(profileJSON{Username: profile.Username, UserID: profile.UserID, Artist: profile.Artist, Location: profile.Location.String, ImageURL: profile.ImageURL.String, HeaderImageURL: profile.HeaderImageURL.String, Description: profile.Description.String, Name: profile.Name}, "Success", http.StatusOK))
}

// ProfileError Renders profile error json
func ProfileError(w http.ResponseWriter, err *models.ProfileError) {
	if err.Internal != "" {
		Error(w, "Internal error, more details in data!", http.StatusNotFound, err)
		return
	}
	Error(w, "Error, more detailts in data", http.StatusBadRequest, err)
}
