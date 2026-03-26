//go:build !PROD

package layout

import (
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AssetEntryPoint(ctx *types.CustomContext) Node {
	entrypoint := "https://assets.hs.localhost/src/main.ts"

	return Script(Type("module"), Src(entrypoint))
}
