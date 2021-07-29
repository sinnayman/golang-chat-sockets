package main

import (
	"fmt"
	"log"
	"net/http"
	"sinnayman/ws/internal/handlers"
	"sync"
)

var port = 8080

func main() {

	var wg sync.WaitGroup

	mux := routes()

	log.Println("Starting web socket listener")
	wg.Add(1)
	go handlers.Listen(wg)

	log.Println(fmt.Sprintf("Starting webserver on port %d", port))
	_ = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)

	wg.Wait()
}
