package phone_test

import (
	"testing"

	"github.com/moroz/homeosapiens-go/internal/phone"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeNumber(t *testing.T) {
	examples := []struct {
		input, expected string
		valid           bool
	}{
		// Different forms of numbers within Poland
		{"+48 555-123-456", "+48555123456", true},
		{"(+48) 555-123-456", "+48555123456", true},
		{"+48555123456", "+48555123456", true},
		{"555-123-456", "+48555123456", true},
		{"+49 171 234 5678", "+491712345678", true},
		{"555123456", "+48555123456", true},
		{"0048555123456", "+48555123456", true},
		{"48555123456", "+48555123456", true},

		// Obviously invalid numbers
		{"", "", false},
		{"abc", "", false},
		{"123", "", false},           // too short
		{"+999555123456", "", false}, // invalid country code
		{"++48555123456", "", false}, // malformed prefix

		// Austria (AT, +43)
		{"+43 688 64159860", "+4368864159860", true},
		// Belgium (BE, +32)
		{"+32 485 49 07 56", "+32485490756", true},
		// Bulgaria (BG, +359)
		{"+359 2 123 4567", "+35921234567", true},
		// Croatia (HR, +385)
		{"+385 1 234 5678", "+38512345678", true},
		// Cyprus (CY, +357)
		{"+357 22 123456", "+35722123456", true},
		// Czech Republic (CZ, +420)
		{"+420 212 345 678", "+420212345678", true},
		// Denmark (DK, +45)
		{"+45 53 77 75 72", "+4553777572", true},
		// Estonia (EE, +372)
		{"+372 512 3456", "+3725123456", true},
		// Finland (FI, +358)
		{"+358 9 1234 5678", "+358912345678", true},
		// France (FR, +33)
		{"+33 1 23 45 67 89", "+33123456789", true},
		// Germany (DE, +49)
		{"+49 30 12345678", "+493012345678", true},
		// Greece (GR, +30)
		{"+30 21 0123 4567", "+302101234567", true},
		// Hungary (HU, +36)
		{"+36 1 234 5678", "+3612345678", true},
		// Ireland (IE, +353)
		{"+353 1 234 5678", "+35312345678", true},
		// Italy (IT, +39)
		{"+39 06 1234 5678", "+390612345678", true},
		// Latvia (LV, +371)
		{"+371-67-483920", "+37167483920", true},
		// Lithuania (LT, +370)
		{"+370 617 80626", "+37061780626", true},
		// Luxembourg (LU, +352)
		{"+352 27 123 456", "+35227123456", true},
		// Malta (MT, +356)
		{"+356 2123 4567", "+35621234567", true},
		// Netherlands (NL, +31)
		{"+31 20 123 4567", "+31201234567", true},
		// Poland (PL, +48)
		{"+48 12 345 6789", "+48123456789", true},
		// Portugal (PT, +351)
		{"+351 21 234 5678", "+351212345678", true},
		// Romania (RO, +40)
		{"+40 21 234 5678", "+40212345678", true},
		// Slovakia (SK, +421)
		{"+421-51-8745621", "+421518745621", true},
		// Slovenia (SI, +386)
		{"+386 1 234 5678", "+38612345678", true},
		// Spain (ES, +34)
		{"+34 91 234 5678", "+34912345678", true},
		// Sweden (SE, +46)
		{"+46 8 123 456 78", "+46812345678", true},
	}

	for _, example := range examples {
		actual, err := phone.NormalizeNumber(example.input)
		if example.valid {
			assert.NoError(t, err)
		}
		assert.Equal(t, example.expected, actual)
	}
}
