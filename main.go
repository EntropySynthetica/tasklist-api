package main

import (
	"fmt"
	"log"
	"net/http"
	"tasklist-api/router"

	"github.com/gorilla/handlers"
)

func main() {
	r := router.Router()

	fmt.Println("Starting server on the port 8080...")

	// log.Fatal(http.ListenAndServe(":8080", r))
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
