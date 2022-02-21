package main

import (
	"log"

	"github.com/teris-io/shortid"
	"golang.org/x/crypto/bcrypt"

	"huangc28/go-ios-iap-vendor/config"
	"huangc28/go-ios-iap-vendor/db"
)

func init() {
	// Init config
	config.InitConfig()

	// Init db
	db.InitDB()
}

func main() {
	username := "admin1"
	password := "1234"

	sid, err := shortid.Generate()
	if err != nil {
		log.Fatal(err)
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt.GenerateFromPassword failed, err: %v", err)
	}

	dbClient := db.GetDB()
	_, err = dbClient.Exec(`
INSERT INTO users (username, password, uuid)
VALUES ($1, $2, $3);
	`, username, hashedPwd, sid)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("seed success %s", hashedPwd)
}
