package countries

import (
	_ "embed"
	"encoding/json"
	"log"
	"slices"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

//go:embed countries.json
var CountryList []byte

type CountryOption struct {
	Value   string
	LabelPl string
	LabelEn string
}

var MostPopularRegionsISO = []string{"PL", "IN", "HK", "CN", "GB", "ZA", "US"}
var EuMemberStatesISO = []string{"AT", "BE", "BG", "HR", "CY", "CZ", "DK", "EE", "FI", "FR", "DE", "GR", "HU", "IE", "IT", "LV", "LT", "LU", "MT", "NL", "PL", "PT", "RO", "SK", "SI", "ES", "SE"}

var all []CountryOption
var mapped map[string]CountryOption
var OrderedByPolish []CountryOption
var OrderedByEnglish []CountryOption
var PopularRegions []CountryOption

func All() []CountryOption {
	return append([]CountryOption(nil), all...)
}

func EuMemberStates() []CountryOption {
	return OptionsFromISOCodeList(EuMemberStatesISO)
}

func SortByLabel(list []CountryOption, langCode string) []CountryOption {
	collation := collate.New(language.English)
	if langCode == "pl" {
		collation = collate.New(language.Polish)
	}

	return slices.SortedFunc(slices.Values(list), func(co1 CountryOption, co2 CountryOption) int {
		switch langCode {
		case "pl":
			return collation.CompareString(co1.LabelPl, co2.LabelPl)
		default:
			return collation.CompareString(co1.LabelEn, co2.LabelEn)
		}
	})
}

func OptionsFromISOCodeList(codes []string) []CountryOption {
	result := make([]CountryOption, 0, len(codes))
	for _, code := range codes {
		if option, ok := mapped[code]; ok {
			result = append(result, option)
		}
	}

	return result
}

func init() {
	if err := json.Unmarshal(CountryList, &all); err != nil {
		log.Fatal(err)
	}
	OrderedByEnglish = SortByLabel(all, "en")
	OrderedByPolish = SortByLabel(all, "pl")

	mapped = make(map[string]CountryOption)
	for i, option := range all {
		mapped[option.Value] = all[i]
	}

	PopularRegions = OptionsFromISOCodeList(MostPopularRegionsISO)
}
