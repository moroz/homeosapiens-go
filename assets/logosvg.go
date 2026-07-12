package svgassets

import _ "embed"

//go:embed public/assets/logo.svg
var LogoSVG string

//go:embed fonts/IBMPlexSans-Bold.ttf
var IBMPlexSansBold []byte

//go:embed fonts/IBMPlexSans-Regular.ttf
var IBMPlexSansRegular []byte
