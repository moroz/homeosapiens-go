package main

import (
	"log"
	"net/http"

	"github.com/moroz/homeosapiens-go/handlers"
)

func main() {
	r := handlers.Router()
	log.Fatal(http.ListenAndServe(":3000", r))
}
