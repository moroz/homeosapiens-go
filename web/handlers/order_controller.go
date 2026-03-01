package handlers

import (
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
)

type orderController struct {
	cartService *services.CartService
}

func OrderController(db queries.DBTX) *orderController {
	return &orderController{
		cartService: services.NewCartService(db),
	}
}
