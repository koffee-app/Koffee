package controllers

import (
	"encoding/json"
	"koffee/internal/middleware"
	"koffee/internal/models"
	"koffee/internal/rabbitmq"
	view "koffee/internal/view"
	"net/http"
	"strconv"
	"strings"

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

// Generalized function for sanitizing body for requests of this controller
func (p *profileBody) sanitize() {
	p.Description = strings.TrimSpace(p.Description)
	p.Username = strings.ToLower(strings.TrimSpace(p.Username))
	p.Name = strings.TrimSpace(p.Name)
}

type profileController struct {
	profileRepo models.RepositoryProfiles
	event       rabbitmq.MessageSender
}

// InitializeProfileController initializes profile routes
func InitializeProfileController(routes *Group, router *httprouter.Router, profile models.RepositoryProfiles, rbEvent rabbitmq.MessageListener) {
	p := profileController{profileRepo: profile, event: rbEvent.NewSender("profile_creation")}
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
	body.sanitize()
	profile, errProfile := p.profileRepo.CreateProfile(body.Username, body.Name, logedUser.UserID, body.Artist)
	if errProfile != nil {
		view.ProfileError(w, errProfile)
		return
	}

	view.Profile(w, profile)
	// Send to the subscribed microservices the event that there is a new profile created
	view.SendJSON(p.event, *profile)
}

// @ PUBLIC GET
// @ PARAMS : identifier string
// @ QUERY : by_username boolean, by_id boolean, is_artist boolean? (optional means that sql won't search by this)
// This getProfile will try to find a profile and how is depending on the params above
// This won't be used for indexed searches!! We'll use another database for that :)
func (p *profileController) getProfile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	queries := r.URL.Query()
	byUsername, byID, isArtist, identifier, useArtist, profile := queries["by_username"], queries["by_id"], queries["is_artist"], params.ByName("identifier"), false, models.Profile{}
	if len(isArtist) > 0 && isArtist[0] != "" {
		useArtist = true
		profile.Artist = isArtist[0] == "true"
	}
	if len(byUsername) > 0 && byUsername[0] == "true" {
		profile.Username = strings.ToLower(identifier)
	} else if len(byID) > 0 && byID[0] == "true" {
		id, err := strconv.Atoi(identifier)
		if err != nil {
			view.ProfileError(w, &models.ProfileError{UserID: "Error parsing byID"})
			return
		}
		profile.UserID = uint32(id)
	}
	profileRef := p.profileRepo.SingleProfile(&profile, useArtist)
	profileRef = p.profileRepo.GetImage(profileRef)
	if profileRef == nil {
		view.ProfileError(w, &models.ProfileError{Internal: "Not found!"})
		return
	}
	view.Profile(w, profileRef)
}

func (p *profileController) updateProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var profile profileBody
	err := json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		view.Error(w, "Error parsing body", http.StatusBadRequest, err)
		return
	}
	artist := strconv.FormatBool(profile.Artist)
	profile.sanitize()
	usr := middleware.GetUser(r)
	profileRes, profileError := p.profileRepo.UpdateProfile(profile.Username, profile.Description, artist, usr.UserID, profile.Name)
	if profileError != nil {
		view.ProfileError(w, profileError)
		return
	}
	view.Profile(w, profileRes)
}

/**
Deleted code that might be useful


if len(multiple) > 0 && multiple[0] != "" {
	multipleParsed, err := strconv.Atoi(multiple[0])
	if err != nil {
		view.ProfileError(w, &models.ProfileError{Internal: "Multiple query param is not a number"})
		return
	}
	profileRef := p.profileRepo.GetProfiles(&profile, useArtist, multipleParsed)
	if profileRef == nil {
		view.ProfileError(w, &models.ProfileError{Internal: "Error retrieving profiles"})
	}
	view.Profiles(w, *profileRef)
	return
}*/
