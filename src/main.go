package main

import (
	"log"
	"net/http"
	"web-example/src/db"
	"web-example/src/handlers"
	"web-example/src/settings"
	"web-example/src/types"

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
