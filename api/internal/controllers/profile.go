package controllers

import (
	"encoding/json"
	"koffee/internal/middleware"
	"koffee/internal/models"
	view "koffee/internal/view"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

type profileBody struct {
	Username    string `json:"username,omitempty"`
	UserID      uint32 `json:"id,omitempty"`
	Artist      bool   `json:"artist,omitempty"`
	Age         uint64 `json:"age,omitempty"`
	ImageURL    string `json:"imageurl,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
}

type profileController struct {
	db *sqlx.DB
}

// InitializeProfileService initializes profile routes
func InitializeProfileService(routes *Group, router *httprouter.Router, db *sqlx.DB) {
	p := profileController{db: db}
	group := New(routes, "/profile")
	router.POST(group.Route("/"), middleware.Authorization(p.createProfile))
	router.GET(group.Route("/:identifier"), p.getProfile)
}

func (p *profileController) createProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logedUser := middleware.GetUser(r)
	body := profileBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		view.ErrorJSON(w, err)
		return
	}
	profile, errProfile := models.CreateProfile(p.db, body.Username, body.Name, logedUser.UserID, body.Artist)
	if errProfile != nil {
		view.ProfileError(w, errProfile)
		return
	}
	view.Profile(w, profile)
}

// @ PUBLIC GET
// @ PARAMS : identifier string
// @ QUERY : by_username boolean, by_id boolean, is_artist boolean? (optional means that sql won't search by this), multiple integer (NOTE: WILL CHANGE THE 200 OK RESPONSE FORMAT OF THIS ROUTE, SEE DOCS)
// This getProfile will try to find a profile and how is depending on the params above
// This won't be used for indexed searches!! We'll use another database for that :)
// TODO: Have an offset for multiple
func (p *profileController) getProfile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	queries := r.URL.Query()
	byUsername, byID, isArtist, identifier, useArtist, profile, multiple := queries["by_username"], queries["by_id"], queries["is_artist"], params.ByName("identifier"), false, models.Profile{}, queries["multiple"]
	if len(isArtist) > 0 && isArtist[0] != "" {
		useArtist = true
		profile.Artist = isArtist[0] == "true"
	}
	if len(byUsername) > 0 && byUsername[0] == "true" {
		profile.Username = identifier
	} else if len(byID) > 0 && byID[0] == "true" {
		id, err := strconv.Atoi(identifier)
		if err != nil {
			view.ProfileError(w, &models.ProfileError{UserID: "Error parsing byID"})
			return
		}
		profile.UserID = uint32(id)
	}
	if len(multiple) > 0 && multiple[0] != "" {
		multipleParsed, err := strconv.Atoi(multiple[0])
		if err != nil {
			view.ProfileError(w, &models.ProfileError{Internal: "Multiple query param is not a number"})
			return
		}
		profileRef := models.GetProfiles(p.db, &profile, useArtist, multipleParsed)
		if profileRef == nil {
			view.ProfileError(w, &models.ProfileError{Internal: "Error retrieving profiles"})
		}
		view.Profiles(w, *profileRef)
		return
	}
	profileRef := profile.GetSingleProfile(p.db, useArtist)
	if profileRef == nil {
		view.ProfileError(w, &models.ProfileError{Internal: "Not found!"})
		return
	}
	view.Profile(w, profileRef)
}
