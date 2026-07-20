package helpers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/types"
)

func GetRequestContext(c *echo.Context) *types.CustomContext {
	return c.Get(config.CustomContextKey).(*types.CustomContext)
}
