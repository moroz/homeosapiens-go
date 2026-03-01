package components

import (
	"github.com/moroz/homeosapiens-go/internal/countries"
	. "maragu.dev/gomponents"
)

type CountrySelectOptions struct {
	Name                  string
	Language              string
	Value                 string
	Required              bool
	Label                 string
	Countries             []countries.CountryOption
	HelperText            string
	IncludePopularRegions bool
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

func buildCountryOptions(list []countries.CountryOption, lang string, includePopular bool) []SelectOption {
	options := countries.SortByLabel(list, lang)

	var combined []SelectOption
	if includePopular {
		combined = mapOptions(countries.PopularRegions, lang)

		combined = append(combined, SelectOption{
			Label: "---",
			Value: "",
		})
	}

	all := mapOptions(options, lang)
	return append(combined, all...)
}

func CountrySelect(opts *CountrySelectOptions) Node {
	options := buildCountryOptions(opts.Countries, opts.Language, opts.IncludePopularRegions)

	name := "country"
	if opts.Name != "" {
		name = opts.Name
	}

	return SelectComponent(&SelectOptions{
		Label:        opts.Label,
		Name:         name,
		Value:        opts.Value,
		Autocomplete: "country",
		Options:      options,
		Required:     opts.Required,
		HelperText:   opts.HelperText,
	})
}
