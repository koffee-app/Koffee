package test

import (
	"database/sql"
	"fmt"
	"koffee/internal/config"
)

var schema = `
CREATE TABLE person (
    first_name text,
    last_name text,
		email text,
		id SERIAL
);

CREATE TABLE place (
    country text,
    city text NULL,
    telcode integer
)`

// Person test struct
type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	ID        uint32 `db:"id"`
}

// Place test struct
type Place struct {
	Country string
	City    sql.NullString
	TelCode int
}

// Posgres test function
func Posgres() {
	db := config.StartConfigDB()
	// This should run if it does not exist.
	// db.MustExec(schema)
	tx := db.MustBegin()
	tx.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Jane", "Citizen", "jane.citzen@example.com", 0})
	tx.Commit()
	people := []Person{}

	db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
	if len(people) == 0 {
		panic("Error in test Posgres line 50")
	}
	fmt.Println(people)
}
