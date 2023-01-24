package main

import (
	"log"

	"github.com/corsairconstantine/http-rest-api/app/apiserver"
)

func main() {
	err := apiserver.Start()
	log.Fatal(err)
}
