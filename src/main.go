package main

import (
	"context"
	"go-api-example/src/db"
	"go-api-example/src/handlers"
	"go-api-example/src/settings"
	"go-api-example/src/types"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	log.Printf("Starting server on %v.", settings.HttpAddr)

	db.InitDB(&db.DB)
	defer db.DB.Close()
	db.DB.AutoMigrate(&types.Wallet{})
	log.Println("Database connected")

	router := handlers.InitRouter()

	server := &http.Server{
		Addr:    settings.HttpAddr,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server critical error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("App has interrupted or cancelled. Start graceful shutdown.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
