package db

import (
	"fmt"
	"log"
)

func main() {
	goquDB := GetConnection()
	exists, err := TableExists(goquDB, "t_fires")
	if err != nil {
		log.Fatalf("Error while checking if table exists: %v", err)
		panic(err)
	}

	if exists {
		fmt.Println("Table t_fires exists")
	} else {
		fmt.Println("Table t_fires doesn't exist, creating it")
		err := CreateTable(goquDB)
		if err != nil {
			log.Fatalf("Error while creating t_fires table: %v", err)
			panic(err)
		}
	}

}