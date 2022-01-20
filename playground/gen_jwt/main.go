package main

import (
	"huangc28/go-ios-iap-vendor/internal/pkg/jwtactor"
	"log"
)

func main() {
	jwt, err := jwtactor.CreateToken("w1YuW907g", "fukyoubitch")

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("jwt %s", jwt)
}
