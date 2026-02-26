package layout

import (
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func navLink(children ...Node) Node {
	return Li(
		Class("-mx-3"),
		A(
			Class("flex h-10 w-full items-center px-3 hover:bg-white/10"),
			Group(children),
		),
	)
}

func AdminLayout(ctx *types.CustomContext, title string, children ...Node) Node {
	return RootLayout(ctx, title,
		Div(
			Class("grid h-screen grid-cols-[300px_1fr]"),
			Aside(
				Class("flex flex-col gap-6 bg-slate-700 p-6 text-white"),
				Header(
					H1(Class("text-2xl font-bold"), Text("Homeo sapiens")),
					H2(Text("Admin panel")),
				),
				Nav(
					Ul(
						Class("grid gap-2"),
						navLink(Href("/admin"), Text("Events")),
						navLink(Href("/admin/users"), Text("Users")),
						navLink(Href("/admin/blog"), Text("Blog")),
						navLink(Href("/"), Text("Back to website")),
					),
				),
				Footer(
					Class("mt-auto"),
					A(Class("button w-full"), Href("/sign-out"), Text("Sign out")),
				),
			),
			Main(
				Class("p-6"),
				Div(Class("container mx-auto max-w-300"), Group(children)),
			),
		),
	)
}
