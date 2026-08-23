package videos

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	svgassets "github.com/moroz/homeosapiens-go/assets"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	. "maragu.dev/gomponents"
)

var logoInnerSVG string
var singleDigitRe = regexp.MustCompile(`\b\d\b`)
var boldFace map[int]font.Face
var regularFace map[int]font.Face
var boldFont *sfnt.Font
var regularFont *sfnt.Font

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

	var err error
	boldFont, err = opentype.Parse(svgassets.DMSansBold)
	if err != nil {
		log.Printf("failed to parse DM Sans Bold: %v", err)
		return
	}
	regularFont, err = opentype.Parse(svgassets.DMSansRegular)
	if err != nil {
		log.Printf("failed to parse DM Sans Regular: %v", err)
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

// textToPath renders text as an SVG path (glyph outlines), so it displays
// correctly regardless of whether the consumer of this standalone SVG has
// the font installed. Baseline is at (x, y) in the same coordinate space as
// font size.
func textToPath(f *sfnt.Font, text string, size, x, y int) string {
	if f == nil {
		return ""
	}
	var buf sfnt.Buffer
	ppem := fixed.I(size)
	penX := fixed.I(x)
	penY := fixed.I(y)

	var sb strings.Builder
	var prev sfnt.GlyphIndex
	for i, r := range text {
		gi, err := f.GlyphIndex(&buf, r)
		if err != nil || gi == 0 {
			continue
		}
		if i > 0 {
			if kern, err := f.Kern(&buf, prev, gi, ppem, font.HintingNone); err == nil {
				penX += kern
			}
		}
		if segs, err := f.LoadGlyph(&buf, gi, ppem, nil); err == nil {
			appendGlyphPath(&sb, segs, penX, penY)
		}
		if adv, err := f.GlyphAdvance(&buf, gi, ppem, font.HintingNone); err == nil {
			penX += adv
		}
		prev = gi
	}
	return sb.String()
}

func appendGlyphPath(sb *strings.Builder, segs sfnt.Segments, offX, offY fixed.Int26_6) {
	for _, seg := range segs {
		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			fmt.Fprintf(sb, "M%s,%s", f26(seg.Args[0].X+offX), f26(seg.Args[0].Y+offY))
		case sfnt.SegmentOpLineTo:
			fmt.Fprintf(sb, "L%s,%s", f26(seg.Args[0].X+offX), f26(seg.Args[0].Y+offY))
		case sfnt.SegmentOpQuadTo:
			fmt.Fprintf(sb, "Q%s,%s %s,%s",
				f26(seg.Args[0].X+offX), f26(seg.Args[0].Y+offY),
				f26(seg.Args[1].X+offX), f26(seg.Args[1].Y+offY))
		case sfnt.SegmentOpCubeTo:
			fmt.Fprintf(sb, "C%s,%s %s,%s %s,%s",
				f26(seg.Args[0].X+offX), f26(seg.Args[0].Y+offY),
				f26(seg.Args[1].X+offX), f26(seg.Args[1].Y+offY),
				f26(seg.Args[2].X+offX), f26(seg.Args[2].Y+offY))
		}
	}
}

func f26(v fixed.Int26_6) string {
	return strconv.FormatFloat(float64(v)/64, 'f', 2, 64)
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

	titlePaths := make([]Node, len(titleLines))
	for i, line := range titleLines {
		lineY := titleLine1Y + i*(titleSize+lineGap)
		titlePaths[i] = El("path", Attr("d", textToPath(boldFont, line, titleSize, contentX, lineY)))
	}

	dateStr := ""
	if date != nil {
		dateStr = formatDate(*date, locale)
	}

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

		El("g",
			Attr("fill", "#0f172a"),
			Group(titlePaths),
		),
		If(subtitle != "",
			El("path",
				Attr("fill", "#334155"),
				Attr("d", textToPath(regularFont, subtitle, subSize, contentX, subtitleY)),
			),
		),

		El("path",
			Attr("fill", "#334155"),
			Attr("d", textToPath(boldFont, host, hostDateSize, contentX, hostY)),
		),
		El("path",
			Attr("fill", "#334155"),
			Attr("d", textToPath(boldFont, dateStr, hostDateSize, contentX, dateY)),
		),
	)
}
