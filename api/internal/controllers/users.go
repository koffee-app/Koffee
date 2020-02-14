package controllers

import (
	"encoding/json"
	"io"
	"koffee/internal/middleware"
	"koffee/internal/models"
	view "koffee/internal/view"
	"koffee/test"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

// user request body
type userBody struct {
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	GoogleToken string `json:"google_access_token,omitempty"`
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
	test.OAUTHGoogle(router)
	router.POST(userGroup.Route("/login"), u.loginHandle)
	router.POST(userGroup.Route("/register"), u.registerHandle)
	router.POST(userGroup.Route("/google"), u.googleAccount)
	router.POST(userGroup.Route("/register/google"), u.googleRegisterHandle)
	router.POST(userGroup.Route("/login/google"), u.googleLoginHandle)
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

func (u *userImpl) googleLoginHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := parseBody(r.Body)
	if err != nil {
		view.Error(w, "Error parsing the request body.", http.StatusBadRequest, err)
		return
	}
	usr, usrerr := models.LogUserGoogle(u.db, body.GoogleToken)
	if usrerr != nil {
		view.RenderAuthError(w, usrerr)
		return
	}
	view.User(w, usr)
}

func (u *userImpl) googleRegisterHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := parseBody(r.Body)

	if err != nil {
		view.Error(w, "Error parsing the request body.", http.StatusBadRequest, err)
		return
	}

	success, notsuccessful, err := models.GoogleAuthentication(body.GoogleToken)

	if err != nil {
		view.Error(w, "Error requesting to Google.", http.StatusNotFound, err)
		return
	}

	if notsuccessful != nil {
		view.Error(w, "Error authenticating with Google access_token", http.StatusBadRequest, notsuccessful)
		return
	}

	if err != nil {
		view.Error(w, "Internal error", http.StatusBadRequest, err)
		return
	}

	usucc, uerr := models.AddUser(u.db, success.Email, "", true)

	if uerr != nil {
		view.RenderAuthError(w, uerr)
		return
	}

	view.User(w, usucc)
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

func parseBody(body io.ReadCloser) (*userBody, error) {
	ub := &userBody{}
	if err := ub.BodyUser(body); err != nil {
		return nil, err
	}
	return ub, nil
}
