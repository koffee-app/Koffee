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

type profilesJSON struct {
	Profiles []profileJSON `json:"profiles"`
}

// profileModelToJSON .
func profileModelToJSON(prof *models.Profile) profileJSON {
	return profileJSON{Username: prof.Username, UserID: prof.UserID, Artist: prof.Artist, Location: prof.Location.String, ImageURL: prof.ImageURL.String, HeaderImageURL: prof.HeaderImageURL.String, Description: prof.Description.String, Name: prof.Name}
}

// Profile Renders profile JSON
func Profile(w http.ResponseWriter, profile *models.Profile) {
	PrepareJSON(w, http.StatusOK)
	ReturnJSON(w, FormatJSON(profileModelToJSON(profile), "Success", http.StatusOK))
}

// Profiles return multiple profiles
func Profiles(w http.ResponseWriter, profilesModel []models.Profile) {
	profiles := profilesJSON{Profiles: make([]profileJSON, len(profilesModel))}
	for i, prof := range profilesModel {
		profiles.Profiles[i] = profileModelToJSON(&prof)
	}
	PrepareJSON(w, http.StatusOK)
	ReturnJSON(w, FormatJSON(profiles, "Success", http.StatusOK))
}

// ProfileError Renders profile error json
func ProfileError(w http.ResponseWriter, err *models.ProfileError) {
	if err.Internal != "" {
		Error(w, "Internal error, more details in data!", http.StatusNotFound, err)
		return
	}
	Error(w, "Error, more details in data", http.StatusBadRequest, err)
}
