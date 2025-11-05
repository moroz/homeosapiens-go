package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type pageController struct{}

func PageController() *pageController {
	return &pageController{}
}

func (c *pageController) Index(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello, world!")
	if err != nil {
		log.Print(err)
	}
}
