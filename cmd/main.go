package main

import (
	"log"

	"github.com/corsairconstantine/http-rest-api/pkg/apiserver"
)

func main() {
	log.Fatal(apiserver.Start())
}
