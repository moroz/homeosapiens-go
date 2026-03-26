package phone

import (
	"errors"
	"regexp"

	"github.com/nyaruka/phonenumbers"
)

var StripRegexp = regexp.MustCompile(`[()\s-]+`)
var ValidationRegexp = regexp.MustCompile(`^\+?[0-9]+$`)

var ValidPhoneNumbers = map[string]string{
	"AT": "+43 688 64159860",
	"BE": "+32 485 49 07 56",
	"BG": "+359 2 123 4567",
	"CY": "+357 22 123456",
	"CZ": "+420 212 345 678",
	"DE": "+49 30 12345678",
	"DK": "+45 53 77 75 72",
	"EE": "+372 512 3456",
	"ES": "+34 91 234 5678",
	"FI": "+358 9 1234 5678",
	"FR": "+33 1 23 45 67 89",
	"GB": "+44 7700 900123",
	"GR": "+30 21 0123 4567",
	"HR": "+385 1 234 5678",
	"HU": "+36 1 234 5678",
	"IE": "+353 1 234 5678",
	"IT": "+39 06 1234 5678",
	"LT": "+370 617 80626",
	"LU": "+352 27 123 456",
	"LV": "+371 67 483920",
	"MT": "+356 2123 4567",
	"NL": "+31 20 123 4567",
	"PL": "+48 12 345 6789",
	"PT": "+351 21 234 5678",
	"RO": "+40 21 234 5678",
	"SE": "+46 8 123 456 78",
	"SI": "+386 1 234 5678",
	"SK": "+421 51 8745621",
}

func NormalizeNumber(p, region string) (string, error) {
	if region == "" {
		region = "PL"
	}

	parsed, err := phonenumbers.Parse(p, "PL")
	if err != nil {
		return "", err
	}

	if !isValidNumber(p) {
		return "", errors.New("invalid number")
	}

	return phonenumbers.Format(parsed, phonenumbers.E164), nil
}

func isValidNumber(p string) bool {
	p = StripRegexp.ReplaceAllString(p, "")

	if !ValidationRegexp.MatchString(p) {
		return false
	}

	parsed, err := phonenumbers.Parse(p, "PL")
	if err != nil {
		return false
	}

	return phonenumbers.IsValidNumber(parsed) && phonenumbers.IsPossibleNumber(parsed)
}

// ExamplePhoneNumber returns a valid phone number for the given country code.
func ExamplePhoneNumber(iso string) string {
	if p, ok := ValidPhoneNumbers[iso]; ok {
		return p
	}
	return ValidPhoneNumbers["PL"]
}
