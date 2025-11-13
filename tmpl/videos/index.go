package videos

import (
	"context"

	"github.com/moroz/homeosapiens-go/tmpl/layout"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Index(ctx context.Context) Node {
	return layout.Layout(ctx, "Videos", H2(Text("Videos")))
}
