package main

import (
	"embox/internal/config"
	"log"

	"embox/internal/infrastructure"
)

func main() {
	dbConfig := config.LoadDbConfig()
	apiConfig := config.LoadApiConfig()

	db, err := infrastructure.InitDatabase(dbConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	infrastructure.InitServer(db, apiConfig)
}
