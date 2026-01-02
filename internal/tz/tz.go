package tz

import (
	_ "embed"
	"strings"
)

//go:embed zone.tab
var tabFile string

var TimezoneMapping map[string]string

type TimezoneGuess struct {
	IsoCode string
	Found   bool
}

func init() {
	TimezoneMapping = make(map[string]string)

	lines := strings.SplitSeq(tabFile, "\n")
	for line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		isoCode := fields[0]
		name := fields[2]

		TimezoneMapping[name] = isoCode
	}
}

func GuessRegionByTimezone(tzName string) *TimezoneGuess {
	iso, ok := TimezoneMapping[tzName]
	return &TimezoneGuess{
		Found:   ok,
		IsoCode: iso,
	}
}
