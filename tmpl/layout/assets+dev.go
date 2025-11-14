//go:build !PROD

package layout

import . "maragu.dev/gomponents"
import . "maragu.dev/gomponents/html"

const IsProd = false

func AssetEntryPoint() Node {
	return Script(Type("module"), Src("http://localhost:5173/src/main.ts"))
}
