package main

import (
	"fmt"
	"log"

	_user "websocket-in-go-boilerplate/src/domains/user"
	_infra_db "websocket-in-go-boilerplate/src/infra/db"
)

func main() {
	runMigrations()
}

func runMigrations() {
	db, err := _infra_db.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("Error on Database Connection: %v", err)
	}

	fmt.Println("Migrating database ...")

	db.AutoMigrate(&_user.User{})

	fmt.Println("Migration completed ...")
}
