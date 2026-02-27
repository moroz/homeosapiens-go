package types

type OrderParams struct {
	BillingGivenName       string
	BillingFamilyName      string
	BillingPhone           *string
	BillingStreet          string
	BillingHouseNumber     *string
	BillingApartmentNumber *string
	BillingCity            string
	BillingPostalCode      *string
	BillingCountry         string
}
