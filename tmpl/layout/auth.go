package layout

import (
	"context"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AuthLayout(ctx context.Context, title string, children ...Node) Node {
	return RootLayout(ctx, title,
		Div(Group(children)))
}
