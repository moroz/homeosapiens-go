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

var MostPopularRegionsISO = []string{"IN", "PL", "HK", "CN", "GB", "ZA", "US"}

var OrderedByPolish []*CountryOption
var OrderedByEnglish []*CountryOption

var PopularRegionsPolish []*CountryOption
var PopularRegionsEnglish []*CountryOption

func init() {
	var options []*CountryOption
	if err := json.Unmarshal(CountryList, &options); err != nil {
		log.Fatal(err)
	}
	colEN := collate.New(language.English)
	colPL := collate.New(language.Polish)
	OrderedByEnglish = slices.SortedFunc(slices.Values(options), func(co1, co2 *CountryOption) int {
		return colEN.CompareString(co1.LabelEn, co2.LabelEn)
	})
	OrderedByPolish = slices.SortedFunc(slices.Values(options), func(co1, co2 *CountryOption) int {
		return colPL.CompareString(co1.LabelPl, co2.LabelPl)
	})

	popularRegions := slices.DeleteFunc(options, func(e *CountryOption) bool {
		return !slices.Contains(MostPopularRegionsISO, e.Value)
	})

	PopularRegionsEnglish = slices.SortedFunc(slices.Values(popularRegions), func(co1 *CountryOption, co2 *CountryOption) int {
		return colEN.CompareString(co1.LabelEn, co2.LabelEn)
	})
	PopularRegionsPolish = slices.SortedFunc(slices.Values(popularRegions), func(co1 *CountryOption, co2 *CountryOption) int {
		return colPL.CompareString(co1.LabelPl, co2.LabelPl)
	})
}
