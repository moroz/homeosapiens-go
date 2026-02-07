package profile

import (
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
)
import . "maragu.dev/gomponents"
import . "maragu.dev/gomponents/html"

func Show(ctx *types.CustomContext) Node {
	return layout.Layout(ctx, "Profile",
		H1(Text("Hello from profile!")),
	)
}
