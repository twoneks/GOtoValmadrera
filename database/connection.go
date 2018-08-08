package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	connStr := "postgres://postgres:@localhost/gotovalma?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	err = db.Ping()
	// Initialize db if does not exit
	if err != nil && err.Error() == "pq: database \"gotovalma\" does not exist" {
		db.Close()
		createDatabase()
	} else if err != nil {
		panic(err)
	}
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func createDatabase() {
	connStr := "postgres://postgres:@localhost/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Initiaizing database...")
	_, err = db.Exec("CREATE DATABASE gotovalma;")
	if err != nil {
		panic(err)
	}
	db.Close()

	// Create table
	connStr = "postgres://postgres:@localhost/gotovalma?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	createTableStatement := `CREATE TABLE detections (
    id          SERIAL NOT NULL PRIMARY KEY,
    knots       integer,
    direction   varchar(3),
    time        TIMESTAMPTZ NOT NULL DEFAULT NOW()
  );
  CREATE TABLE stats (
    id          SERIAL NOT NULL PRIMARY KEY,
    averages    json,
    day         date,
    windy       boolean
  );
  `

	_, err = db.Exec(createTableStatement)

	if err != nil {
		panic(err)
	}

	fmt.Println("Database ready!")
}
