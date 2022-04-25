package main

import (
	"huangc28/go-ios-iap-vendor/internal/pkg/jwtactor"
	"log"
)

func main() {
	jwt, err := jwtactor.CreateToken("t-aOdsb7g", "fukyoubitch")

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("jwt %s", jwt)
}
