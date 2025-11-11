package layout

import (
	"context"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AuthLayout(ctx context.Context, title string, children ...Node) Node {
	return RootLayout(ctx, title,
		Div(
			Class("bg-slate-100 grid place-items-center h-screen gap-4"),
			Div(
				Class("flex flex-col gap-4 items-center -mt-25"),
				A(Href("/"), Class("text-primary hover:text-primary-hover transition-colors"),
					logo("h-15"),
				),
				Div(Class("card w-100"),
					H2(
						Class("text-3xl text-primary font-semibold text-center"),
						Text(title),
					),
					Group(children),
				),
			),
		),
	)
}
