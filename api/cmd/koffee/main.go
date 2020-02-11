package main

import (
	"fmt"
	"koffee/internal/config"
	"koffee/internal/controllers"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func startServer(db *sqlx.DB) {
	router := httprouter.New()
	group := controllers.Group{Prefix: "/api"}
	controllers.InitializeUserController(&group, router, db)
	fmt.Println("Connected on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	startServer(config.InitConfig())
}
