package main

import (
	"fmt"
	"github.com/vadimlarionov/expert-system/es"
	"github.com/vadimlarionov/expert-system/model"
)

func main() {
	username := "es_user"
	password := "es_password"
	dbName := "es_db"

	err := model.Init(username, password, dbName, true)
	if err != nil {
		fmt.Printf("Can't init models: %s", err)
	}

	es.Fill()
}
