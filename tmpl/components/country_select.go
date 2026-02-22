package components

import (
	"github.com/moroz/homeosapiens-go/internal/countries"
	. "maragu.dev/gomponents"
)

type CountrySelectOptions struct {
	Language string
	Value    string
	Required bool
	Label    string
}

func mapOptions(options []countries.CountryOption, lang string) []SelectOption {
	var combined []SelectOption
	for _, o := range options {
		label := o.LabelEn
		if lang == "pl" {
			label = o.LabelPl
		}

		combined = append(combined, SelectOption{
			Label: label,
			Value: o.Value,
		})
	}
	return combined
}

func buildCountryOptions(lang string) []SelectOption {
	options := countries.SortByLabel(countries.All(), lang)

	var combined []SelectOption
	combined = mapOptions(countries.PopularRegions, lang)

	combined = append(combined, SelectOption{
		Label: "---",
		Value: "",
	})

	all := mapOptions(options, lang)
	return append(combined, all...)
}

func CountrySelect(opts *CountrySelectOptions) Node {
	options := buildCountryOptions(opts.Language)

	return SelectComponent(&SelectOptions{
		Label:        opts.Label,
		Name:         "country",
		Value:        opts.Value,
		Autocomplete: "country",
		Options:      options,
		Required:     opts.Required,
	})
}
