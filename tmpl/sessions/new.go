package sessions

import (
	"context"

	"github.com/moroz/homeosapiens-go/tmpl/layout"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func New(ctx context.Context) Node {
	return layout.AuthLayout(ctx, "Sign in",
		Div(
			Class("card"),
			H3(Text("Sign in")),
		),
	)
}
