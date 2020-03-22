package view

import (
	"koffee/internal/models"
	"net/http"
)

type profileJSON struct {
	Username    string `json:"username,omitempty"`
	UserID      uint32 `json:"id,omitempty"`
	Artist      bool   `json:"artist"`
	Age         uint64 `json:"age,omitempty"`
	ImageURL    string `json:"imageurl,omitempty"`
	Description string `json:"description,omitempty"`
}

// Profile Renders profile JSON
func Profile(w http.ResponseWriter, profile *models.Profile) {
	age, _ := profile.Age.Value()
	var ageUint uint64
	if age == nil {
		ageUint = 0
	} else {
		ageUint = age.(uint64)
	}
	PrepareJSON(w, http.StatusOK)
	ReturnJSON(w, FormatJSON(profileJSON{Username: profile.Username, UserID: profile.UserID, Artist: profile.Artist, Age: ageUint, ImageURL: profile.ImageURL.String, Description: profile.Description.String}, "Success", http.StatusOK))
}

// ProfileError Renders profile error json
func ProfileError(w http.ResponseWriter, err *models.ProfileError) {
	if err.Internal != "" {
		Error(w, "Internal error, more details in data!", http.StatusNotFound, err)
		return
	}
	Error(w, "Error, more detailts in data", http.StatusBadRequest, err)
}
