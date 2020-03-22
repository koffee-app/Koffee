package controllers

import (
	"encoding/json"
	"koffee/internal/middleware"
	"koffee/internal/models"
	view "koffee/internal/view"
	"net/http"

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
}

type profileController struct {
	db *sqlx.DB
}

// InitializeProfileService initializes profile routes
func InitializeProfileService(routes *Group, router *httprouter.Router, db *sqlx.DB) {
	p := profileController{db: db}
	group := New(routes, "/profile")
	router.POST(group.Route("/"), middleware.Authorization(p.createProfile))
}

func (p *profileController) createProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logedUser := middleware.GetUser(r)
	body := profileBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		view.ErrorJSON(w, err)
		return
	}
	profile, errProfile := models.CreateProfile(p.db, body.Username, logedUser.UserID, body.Artist)
	if errProfile != nil {
		view.ProfileError(w, errProfile)
		return
	}
	view.Profile(w, profile)
}
