package config

import (
	"fmt"
	"koffee/pkg/countries"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// InitConfig inits the configuration files and more
func InitConfig() *sqlx.DB {
	countries.Init()
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
		panic(".env file not found.")
	}
	JWTConfig()
	PasetoInit()
	return StartConfigDB()
}
