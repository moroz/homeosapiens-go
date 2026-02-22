package orders

import (
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
)
import (
	. "maragu.dev/gomponents"
)

func New(ctx *types.CustomContext) Node {
	return layout.Layout(ctx, "Checkout")
}
