package settings

import "os"

var HttpAddr = os.Getenv("HTTP_ADDR")

var DBHost = os.Getenv("DB_HOST")
var DBPort = os.Getenv("DB_PORT")
var DBUser = os.Getenv("POSTGRES_USER")
var DBPass = os.Getenv("POSTGRES_PASSWORD")
var DBName = os.Getenv("POSTGRES_DB")
