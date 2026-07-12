package videos

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	svgassets "github.com/moroz/homeosapiens-go/assets"
	. "maragu.dev/gomponents"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var logoInnerSVG string
var singleDigitRe = regexp.MustCompile(`\b\d\b`)
var boldFace map[int]font.Face
var regularFace map[int]font.Face

func init() {
	content := svgassets.LogoSVG
	svgStart := strings.Index(content, "<svg")
	if svgStart < 0 {
		return
	}
	tagEnd := strings.Index(content[svgStart:], ">") + svgStart + 1
	svgEnd := strings.LastIndex(content, "</svg>")
	if tagEnd > 0 && svgEnd > tagEnd {
		logoInnerSVG = content[tagEnd:svgEnd]
	}

	boldFont, err := opentype.Parse(svgassets.IBMPlexSansBold)
	if err != nil {
		log.Printf("failed to parse IBM Plex Sans Bold: %v", err)
		return
	}
	regularFont, err := opentype.Parse(svgassets.IBMPlexSansRegular)
	if err != nil {
		log.Printf("failed to parse IBM Plex Sans Regular: %v", err)
		return
	}

	boldFace = make(map[int]font.Face)
	regularFace = make(map[int]font.Face)
	for _, size := range []int{24, 28, 30, 36, 40} {
		bf, err := opentype.NewFace(boldFont, &opentype.FaceOptions{Size: float64(size), DPI: 72})
		if err == nil {
			boldFace[size] = bf
		}
		rf, err := opentype.NewFace(regularFont, &opentype.FaceOptions{Size: float64(size), DPI: 72})
		if err == nil {
			regularFace[size] = rf
		}
	}
}

func measureText(f font.Face, text string) int {
	if f == nil {
		return len([]rune(text)) * 20
	}
	advance := font.MeasureString(f, text)
	return int(advance >> 6)
}

func splitTitle(baseTitle string) (title, subtitle string) {
	parts := strings.SplitN(baseTitle, ":", 2)
	if len(parts) == 2 {
		sub := strings.TrimSpace(parts[1])
		if singleDigitRe.MatchString(sub) {
			return strings.TrimSpace(parts[0]), sub
		}
	}
	return baseTitle, ""
}

func formatDate(t time.Time, locale string) string {
	if locale == "pl" {
		months := [13]string{"", "stycznia", "lutego", "marca", "kwietnia", "maja", "czerwca",
			"lipca", "sierpnia", "września", "października", "listopada", "grudnia"}
		return fmt.Sprintf("%d %s %d", t.Day(), months[t.Month()], t.Year())
	}
	return t.Format("January 2, 2006")
}

func px(n int) string { return strconv.Itoa(n) }

// wrapTextFace splits text into at most maxLines lines using actual font metrics.
func wrapTextFace(f font.Face, text string, maxWidth, maxLines int) []string {
	return wrapWordsFace(f, strings.Fields(text), maxWidth, maxLines)
}

func wrapWordsFace(f font.Face, words []string, maxWidth, maxLines int) []string {
	if len(words) == 0 {
		return nil
	}
	if maxLines <= 1 {
		return []string{strings.Join(words, " ")}
	}
	line := ""
	for i, word := range words {
		candidate := line
		if candidate != "" {
			candidate += " "
		}
		candidate += word
		if measureText(f, candidate) > maxWidth {
			if line == "" {
				return append([]string{word}, wrapWordsFace(f, words[i+1:], maxWidth, maxLines-1)...)
			}
			return append([]string{line}, wrapWordsFace(f, words[i:], maxWidth, maxLines-1)...)
		}
		line = candidate
	}
	return []string{line}
}

func Thumbnail(baseTitle, host, locale string, date *time.Time, ppURL *string) Node {
	title, subtitle := splitTitle(baseTitle)
	smallerText := subtitle != "" || len([]rune(baseTitle)) > 33

	logoW, logoH := 220, 51
	titleSize, subSize, hostDateSize := 40, 30, 30
	if smallerText {
		logoW, logoH = 200, 47
		titleSize, subSize, hostDateSize = 36, 28, 24
	}

	const contentX = 280
	const contentW = 328 // 640 - 32(right pad) - contentX
	const topY = 24
	const bottomY = 336

	titleLines := wrapTextFace(boldFace[titleSize], title, contentW, 3)

	// Block heights (in SVG units)
	const lineGap = 2
	titleBlockH := len(titleLines)*titleSize + (len(titleLines)-1)*lineGap
	if subtitle != "" {
		titleBlockH += 8 + subSize
	}
	hostDateBlockH := hostDateSize*2 + 4

	// justify-between: 3 blocks evenly spaced within topY..bottomY
	totalH := logoH + titleBlockH + hostDateBlockH
	gap := (bottomY - topY - totalH) / 2

	logoY := topY
	titleTop := logoY + logoH + gap
	hostTop := titleTop + titleBlockH + gap

	titleLine1Y := titleTop + titleSize
	lastTitleLineY := titleLine1Y + (len(titleLines)-1)*(titleSize+lineGap)

	var subtitleY int
	if subtitle != "" {
		subtitleY = lastTitleLineY + 8 + subSize
	}

	hostY := hostTop + hostDateSize
	dateY := hostY + 4 + hostDateSize

	tspans := make([]Node, len(titleLines))
	for i, line := range titleLines {
		attrs := []Node{Attr("x", px(contentX))}
		if i > 0 {
			attrs = append(attrs, Attr("dy", px(titleSize+lineGap)))
		}
		tspans[i] = El("tspan", append(attrs, Text(line))...)
	}

	dateStr := ""
	if date != nil {
		dateStr = formatDate(*date, locale)
	}

	const fontFamily = `"IBM Plex Sans", system-ui, sans-serif`

	return El("svg",
		Attr("xmlns", "http://www.w3.org/2000/svg"),
		Attr("viewBox", "0 0 640 360"),
		Attr("width", "1280"),
		Attr("height", "720"),

		El("rect", Attr("width", "640"), Attr("height", "360"), Attr("fill", "#f1f5f9")),

		El("defs",
			El("clipPath", Attr("id", "photo-clip"),
				El("rect", Attr("x", "16"), Attr("y", "16"), Attr("width", "240"), Attr("height", "328"), Attr("rx", "6")),
			),
		),

		Iff(ppURL != nil, func() Node {
			return El("image",
				Attr("x", "16"), Attr("y", "16"),
				Attr("width", "240"), Attr("height", "328"),
				Attr("href", *ppURL),
				Attr("preserveAspectRatio", "xMidYMid slice"),
				Attr("clip-path", "url(#photo-clip)"),
			)
		}),
		If(ppURL == nil,
			El("rect", Attr("x", "16"), Attr("y", "16"), Attr("width", "240"), Attr("height", "328"), Attr("rx", "6"), Attr("fill", "#e2e8f0")),
		),
		El("rect",
			Attr("x", "16"), Attr("y", "16"), Attr("width", "240"), Attr("height", "328"),
			Attr("rx", "6"), Attr("fill", "none"), Attr("stroke", "#64748b"), Attr("stroke-width", "1"),
		),

		El("svg",
			Attr("x", px(contentX)), Attr("y", px(logoY)),
			Attr("width", px(logoW)), Attr("height", px(logoH)),
			Attr("viewBox", "0 0 1538 361"),
			Attr("color", "#0f172a"),
			Raw(logoInnerSVG),
		),

		El("text",
			Attr("x", px(contentX)), Attr("y", px(titleLine1Y)),
			Attr("font-size", px(titleSize)), Attr("font-weight", "bold"),
			Attr("fill", "#0f172a"), Attr("font-family", fontFamily),
			Group(tspans),
		),
		If(subtitle != "",
			El("text",
				Attr("x", px(contentX)), Attr("y", px(subtitleY)),
				Attr("font-size", px(subSize)),
				Attr("fill", "#334155"), Attr("font-family", fontFamily),
				Text(subtitle),
			),
		),

		El("text",
			Attr("x", px(contentX)), Attr("y", px(hostY)),
			Attr("font-size", px(hostDateSize)), Attr("font-weight", "600"),
			Attr("fill", "#334155"), Attr("font-family", fontFamily),
			Text(host),
		),
		El("text",
			Attr("x", px(contentX)), Attr("y", px(dateY)),
			Attr("font-size", px(hostDateSize)), Attr("font-weight", "600"),
			Attr("fill", "#334155"), Attr("font-family", fontFamily),
			Text(dateStr),
		),
	)
}
