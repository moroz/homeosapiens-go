package videos

import (
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Youtube(ctx *types.CustomContext, videos []*queries.Video) Node {
	return layout.BareLayout(ctx, "Watch", Div())
}
