package main

import (
	"go-api-example/src/db"
	"go-api-example/src/handlers"
	"go-api-example/src/settings"
	"go-api-example/src/types"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	log.Printf("Starting web-server on %v.", settings.HttpAddr)

	db.InitDB(&db.DB)
	defer db.DB.Close()
	db.DB.AutoMigrate(&types.Wallet{})
	log.Println("Database connected")

	router := handlers.InitRouter()

	server := &http.Server{
		Addr:    settings.HttpAddr,
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Web-server critical error: %v", err)
	}
}
