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
			Class("desktop:grid desktop:place-items-center min-h-screen gap-4 bg-slate-100"),
			Div(
				Class("desktop:-mt-25 flex w-full flex-col items-center gap-4"),
				A(Href("/"), Class("text-primary hover:text-primary-hover mobile:mt-8 transition-colors"),
					logo("h-15"),
				),
				Div(Class("card desktop:w-100 w-full"),
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
