package db

import (
	"log"
	"web-example/src/settings"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var DB *gorm.DB

func InitDB(db **gorm.DB) {
	var err error

	*db, err = gorm.Open(
		"postgres",
		"host="+settings.DBHost+" port="+settings.DBPort+" user="+settings.DBUser+
			" dbname="+settings.DBName+" sslmode=disable password="+settings.DBPass)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}
}
