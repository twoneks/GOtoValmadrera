package main

import (
	"fmt"

	"github.com/twoneks/gotovalma/database"
)

func main() {
	fmt.Println("Booting...")
	db := database.Connect()
	_, err := db.Exec("INSERT INTO wind VALUES (DEFAULT, 10, 'NE')")
	if err != nil {
		panic(err)
	}

	defer db.Close()
}
