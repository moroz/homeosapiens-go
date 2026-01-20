package dashboard

import (
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
)

func Index(ctx *types.CustomContext) Node {
	return layout.Layout(ctx, "Dashboard", Text("Dashboard"))
}
