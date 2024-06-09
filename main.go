package main

import (
	"log"

	"github.com/mxnyawi/gymSharkTask/internal/db"
	"github.com/mxnyawi/gymSharkTask/pkg/api"
)

func main() {
	dbManager, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	api.StartServer(dbManager)
}
