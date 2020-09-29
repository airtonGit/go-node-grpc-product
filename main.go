package main

import (
	"log"

	"github.com/airtonGit/go-node-grpc-product/rest"
)

func main() {
	log.Printf("Server started")

	log.Fatal(rest.Listen())
}
