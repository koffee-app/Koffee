package main

import (
	"fmt"
	"koffee/internal/config"
	"koffee/internal/controllers"
	"koffee/internal/models"
	"koffee/internal/view"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func startServer(db *sqlx.DB) {
	// initialize tables
	models.Initialize(db)
	// initialize controllers
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		view.PrepareJSON(w, http.StatusOK)
		view.ReturnJSON(w, view.FormatJSON(nil, "Test succeeded! You can use Koffee!", 200))
	})
	group := controllers.Group{Prefix: "/api"}
	controllers.InitializeUserController(&group, router, db)
	controllers.InitializeDriverController(&group, router, db)
	controllers.InitializeProfileService(&group, router, db)
	// Inform that we finished intiializing
	fmt.Println("Connected on port 8081")
	if http.ListenAndServe(":8081", router) != nil {
		fmt.Println("error")
	}
}

func main() {
	startServer(config.InitConfig())
}
