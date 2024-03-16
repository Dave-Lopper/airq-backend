package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/lib/pq"
)

var dbUserName string = os.Getenv("DB_USERNAME")
var dbPassword string = os.Getenv("DB_PASSWORD")
var dbName string = os.Getenv("DB_NAME")
var connStr string = fmt.Sprintf("user=%v dbname=%v sslmode=disable", dbUserName, dbName)


func GetConnection() *goqu.Database {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error while connecting DB: %v", err)
	}
	return goqu.New("postgres", db)
}