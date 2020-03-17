package controllers

import (
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

type profileBody struct {
	Username    string `db:"username"`
	UserID      uint32 `db:"id"`
	Artist      bool   `db:"artist"`
	Age         uint64 `db:"age"`
	ImageURL    string `db:"imageurl"`
	Description string `db:"description"`
}

type profileController struct {
	db *sqlx.DB
}

// InitializeProfileService initializes profile routes
func InitializeProfileService(routes *Group, router *httprouter.Router, db *sqlx.DB) {
	// p := profileController{db: db}
	// group := New(routes, "/profile")
	// router.POST(group.Route("/"), middleware.Authorization())
}
