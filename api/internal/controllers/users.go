package controllers

import (
	"encoding/json"
	"io"
	"koffee/internal/middleware"
	"koffee/internal/models"
	view "koffee/internal/view"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

// user request body
type userBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *userBody) BodyUser(r io.ReadCloser) error {
	return json.NewDecoder(r).Decode(u)
}

// user controller implementation
type userImpl struct {
	db *sqlx.DB
}

// InitializeUserController inits the routes for the routes.
func InitializeUserController(api *Group, router *httprouter.Router, db *sqlx.DB) {
	u := userImpl{db: db}
	userGroup := New(api, "/user")
	models.Initialize(u.db)
	router.POST(userGroup.Route("/login"), u.loginHandle)
	router.POST(userGroup.Route("/register"), u.registerHandle)
	router.POST(userGroup.Route("/google"), u.googleAccount)
	router.GET(userGroup.Route(""), middleware.JwtAuthentication(u.currentUser))
}

func (u *userImpl) currentUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view.User(w, middleware.GetUser(r))
}

func (u *userImpl) loginHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body := &userBody{}
	if err := body.BodyUser(r.Body); err != nil {
		view.Error(w, "Error parsing the request body.", http.StatusBadRequest, err)
		return
	}
	if usuccess, uerr := models.LogUserNotGoogle(u.db, body.Email, body.Password); uerr != nil {
		view.RenderAuthError(w, uerr)
	} else {
		view.User(w, usuccess)
	}
}

// handles the registration route
func (u *userImpl) registerHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body := &userBody{}
	if err := body.BodyUser(r.Body); err != nil {
		view.Error(w, "Error parsing the request body.", http.StatusBadRequest, err)
		return
	}
	if usuccess, uerr := models.AddUser(u.db, body.Email, body.Password, false); uerr != nil {
		view.RenderAuthError(w, uerr)
	} else {
		view.User(w, usuccess)
	}
}

// Verify or create account after signin with google
func (u *userImpl) googleAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}
