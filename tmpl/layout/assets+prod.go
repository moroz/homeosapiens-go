//go:build PROD

package layout

import . "maragu.dev/gomponents"

const IsProd = true

func AssetEntryPoint() Node {
	return Group{}
}
