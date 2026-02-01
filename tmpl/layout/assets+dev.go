//go:build !PROD

package layout

import (
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

const IsProd = false

func AssetEntryPoint(ctx *types.CustomContext) Node {
	entrypoint := "http://" + ctx.RequestUrl.Hostname() + ":5173/src/main.ts"

	return Script(Type("module"), Src(entrypoint))
}
