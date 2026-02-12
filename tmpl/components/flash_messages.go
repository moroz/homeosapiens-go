package components

import (
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Flash(messages types.Flash) Node {
	var elements []Node

	for level, msg := range messages {
		alert := Article(Class("alert "+level), Text(msg))
		elements = append(elements, alert)
	}

	return Group(elements)
}
