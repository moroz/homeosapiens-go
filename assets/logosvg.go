package svgassets

import _ "embed"

//go:embed public/assets/logo.svg
var LogoSVG string

//go:embed fonts/DMSans-Bold.ttf
var DMSansBold []byte

//go:embed fonts/DMSans-Regular.ttf
var DMSansRegular []byte
