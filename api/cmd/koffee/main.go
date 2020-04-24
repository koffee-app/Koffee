package main

import (
	"fmt"
	"koffee/internal/auth"
	"koffee/internal/config"
	"koffee/internal/controllers"
	"koffee/internal/middleware"
	"koffee/internal/models"
	"koffee/internal/rabbitmq"
	"koffee/internal/view"
	"koffee/test/repository"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func startServer(db *sqlx.DB) {
	rabbit := rabbitmq.Initialize()
	tokenService := auth.NewPaseto()

	middleware.Initialize(tokenService)
	// initialize repos
	users, profiles, albums, _ := models.Initialize(db, tokenService)
	// TODO Put this in a single function that tests everything we want
	repository.TestCreateAlbum(profiles, albums, db)
	// initialize controllers
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		view.PrepareJSON(w, http.StatusOK)
		view.ReturnJSON(w, view.JSON(nil, "Test succeeded! You can use Koffee!", 200))
	})
	group := controllers.Group{Prefix: "/api"}
	controllers.InitializeUserController(&group, router, users)
	controllers.InitializeProfileController(&group, router, profiles, rabbit)
	controllers.InitializeAlbumController(&group, router, albums, rabbit)
	// Inform that we finished intiializing
	fmt.Println("Connected on port 8080")
	if http.ListenAndServe(":8080", middleware.ApplyCors(router, &middleware.Cors{Origin: "*"})) != nil {
		fmt.Println("error")
	}
}

func main() {
	startServer(config.InitConfig())
}
