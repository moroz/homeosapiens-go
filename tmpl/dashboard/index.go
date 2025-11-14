package dashboard

import (
	"context"

	"github.com/moroz/homeosapiens-go/tmpl/layout"
	. "maragu.dev/gomponents"
)

func Index(ctx context.Context) Node {
	return layout.Layout(ctx, "Dashboard", Text("Dashboard"))
}
