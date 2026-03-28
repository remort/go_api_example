package db

import (
	"go-api-example/src/settings"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var DB *gorm.DB

func InitDB(db **gorm.DB) {
	var err error
	maxRetries := 10
	retryInterval := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		*db, err = gorm.Open("postgres", getDSN())
		if err == nil {
			if err = (*db).DB().Ping(); err == nil {
				log.Println("Successfully connected to database")
				return
			}
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %s", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	log.Fatalf("Failed to connect to database after %d attempts: %s", maxRetries, err)
}

func getDSN() string {
	return "host=" + settings.DBHost +
		" port=" + settings.DBPort +
		" user=" + settings.DBUser +
		" dbname=" + settings.DBName +
		" sslmode=disable password=" + settings.DBPass
}
