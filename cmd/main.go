package main

import (
	"go-rinha-de-backend-2023/config"
	"log"
)

func main() {
	err := config.InitializeRouter()

	if err != nil {
		log.Fatalln("failed to initialize router with error: ", err)
	}
}
