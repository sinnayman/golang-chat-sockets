package main

import (
	"fmt"
	"log"
	"net/http"
)

var port = 8080

func main() {

	mux := routes()

	log.Println(fmt.Sprintf("Starting webserver on port %d", port))
	_ = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
