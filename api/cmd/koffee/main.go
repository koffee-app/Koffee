package main

import (
	"fmt"
	"koffee/internal/config"
	"koffee/internal/controllers"
	"koffee/internal/view"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func startServer(db *sqlx.DB) {
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		view.ReturnJSON(w, view.FormatJSON(nil, "Test succeeded! You can use Koffee!.", 200))
	})
	group := controllers.Group{Prefix: "/api"}
	controllers.InitializeUserController(&group, router, db)
	fmt.Println("Connected on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	startServer(config.InitConfig())
}
