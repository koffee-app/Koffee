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

// The driver controller request body
type driverBody struct {
	Fullname string `json:"fullname,omitempty"`
	UserID   uint32 `json:"id,omitempty"`
	Country  string `json:"country,omitempty"`
	// (GABI): Consideration, maybe more will be included here
	//				 like finding the current driver of a passenger or sth like that
}

type driverImpl struct {
	db *sqlx.DB
}

// InitializeDriverController .
func InitializeDriverController(api *Group, router *httprouter.Router, db *sqlx.DB) {
	d := driverImpl{db: db}
	driverGroup := New(api, "/driver")
	router.POST(driverGroup.Route("/"), middleware.JwtAuthentication(d.createDriver))
}

// createDriver creates a new driver handler
func (d *driverImpl) createDriver(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := middleware.GetUser(r)
	dbody := driverBody{}

	if err := json.NewDecoder(r.Body).Decode(&dbody); err != nil {
		view.Error(w, "Error parsing the request body.", http.StatusBadRequest, err)
		return
	}
	// Set userID because we don't receive it from the request in authentication
	dbody.UserID = u.UserID
	if driverAlreadyExists := models.GetDriverByID(d.db, dbody.UserID); driverAlreadyExists != nil {
		view.Error(w, "Driver already exists!", http.StatusNotFound, &view.DriverJSON{ImageURL: "", Fullname: driverAlreadyExists.Fullname, Country: driverAlreadyExists.Country, UserID: driverAlreadyExists.UserID})
		return
	}
	dsucc, derr := models.CreateDriver(d.db, dbody.UserID, dbody.Fullname, dbody.Country)
	if derr != nil {
		view.Error(w, "Error creating driver.", http.StatusBadRequest, derr)
		return
	}
	view.Driver(w, dsucc)
}
