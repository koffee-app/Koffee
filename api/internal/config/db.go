package config

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // here
)

func returnConfigSQLString() string {
	userDB, _ := os.LookupEnv("DB_USER")
	nameDB, _ := os.LookupEnv("DB_NAME")
	password, _ := os.LookupEnv("DB_PASSWORD")
	hostDB, _ := os.LookupEnv("DB_HOST")
	portDB, _ := os.LookupEnv("DB_PORT")
	s := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", hostDB, portDB, userDB, nameDB, password)
	fmt.Println(s)
	return s
}

// StartConfigDB configs the Postgres database
func StartConfigDB() *sqlx.DB {

	var err error

	db, err := sqlx.Connect("postgres", returnConfigSQLString())

	for err != nil {
		fmt.Println("Trying...")
		time.Sleep(2 * time.Second)
		db, err = sqlx.Connect("postgres", returnConfigSQLString())
		fmt.Println(err)
	}

	fmt.Println("POSTGRES CONNECTED!")

	return db
}
