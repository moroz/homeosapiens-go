package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func SignOutButton(children ...Node) Node {
	return Form(
		Action("/sign-out"),
		Method("POST"),
		Input(
			Type("hidden"),
			Name("_method"),
			Value("DELETE"),
		),
		Button(
			Type("submit"),
			Group(children),
		),
	)
}
