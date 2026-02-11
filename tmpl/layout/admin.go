package layout

import (
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AdminLayout(ctx *types.CustomContext, title string, children ...Node) Node {
	return RootLayout(ctx, title,
		Div(
			Class("grid h-screen grid-cols-[300px_1fr]"),
			Aside(
				Class("bg-[salmon]"),
				H1(Text("Admin panel")),
			),
			Main(
				Group(children),
			),
		),
	)
}
