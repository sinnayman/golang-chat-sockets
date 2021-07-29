package main

import (
	"fmt"
	"log"
	"net/http"
	"sinnayman/ws/internal/handlers"
)

var port = 8080

func main() {
	mux := routes()

	log.Println("Starting web socket listener")
	go handlers.Listen()

	log.Println(fmt.Sprintf("Starting webserver on port %d", port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Println(err)
	}

}
