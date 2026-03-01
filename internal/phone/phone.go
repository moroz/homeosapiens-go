package phone

import (
	"errors"
	"regexp"

	"github.com/nyaruka/phonenumbers"
)

var StripRegexp = regexp.MustCompile(`[()\s-]+`)
var ValidationRegexp = regexp.MustCompile(`^\+?[0-9]+$`)

func NormalizeNumber(p string) (string, error) {
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
