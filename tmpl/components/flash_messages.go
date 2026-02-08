package components

import (
	"strings"

	"github.com/moroz/homeosapiens-go/types"
)
import . "maragu.dev/gomponents"
import . "maragu.dev/gomponents/html"

func Flash(messages []types.FlashMessage) Node {
	return Map(messages, func(msg types.FlashMessage) Node {
		return Article(
			Class("alert "+strings.ToLower(msg.Level.String())),
			Text(msg.Message),
		)
	})
}
