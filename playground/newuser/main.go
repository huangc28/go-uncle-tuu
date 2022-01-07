package main

import (
	"huangc28/go-ios-iap-vendor/config"
	"huangc28/go-ios-iap-vendor/db"
	"log"

	"github.com/teris-io/shortid"
)

func init() {
	// Init config
	config.InitConfig()

	// Init db
	db.InitDB()
}

func main() {
	username := "admin"
	password := "1234"
	sid, err := shortid.Generate()

	if err != nil {
		log.Fatal(err)
	}

	dbClient := db.GetDB()
	_, err = dbClient.Exec(`
INSERT INTO users (username, password, uuid)
VALUES ($1, $2, $3);
`, username, password, sid)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("seed success")
}
