package main

import (
	"Server/server"
	"log"
)

func main() {
	if err := server.Start(); err != nil {
		log.Println(err)
	}
}
