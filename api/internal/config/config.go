package config

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// InitConfig inits the configuration files and more
func InitConfig() *sqlx.DB {

	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
		panic(".env file not found.")
	}
	JWTConfig()
	return StartConfigDB()
}
