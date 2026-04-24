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
			Class("min-h-screen gap-4 bg-slate-100 desktop:grid desktop:place-items-center"),
			Div(
				Class("flex w-full flex-col items-center gap-4 desktop:-mt-25"),
				A(Href("/"), Class("text-primary transition-colors hover:text-primary-hover mobile:mt-8"),
					components.Logo("h-15"),
				),
				Div(Class("card w-full desktop:w-100"),
					H2(
						Class("text-center text-3xl font-semibold text-primary"),
						Text(title),
					),
					components.Flash(ctx.Flash),
					Group(children),
				),
			),
		),
	)
}
