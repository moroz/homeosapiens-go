package layout

import (
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AuthLayout(ctx *types.CustomContext, title string, children ...Node) Node {
	return RootLayout(ctx, title,
		Div(
			Class("grid h-screen place-items-center gap-4 bg-slate-100"),
			Div(
				Class("-mt-25 flex flex-col items-center gap-4"),
				A(Href("/"), Class("text-primary hover:text-primary-hover transition-colors"),
					logo("h-15"),
				),
				Div(Class("card w-100"),
					H2(
						Class("text-primary text-center text-3xl font-semibold"),
						Text(title),
					),
					components.Flash(ctx.Flash),
					Group(children),
				),
			),
		),
	)
}
