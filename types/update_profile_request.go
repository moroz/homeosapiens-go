package types

type UpdateProfileRequest struct {
	GivenName     string `form:"given_name"`
	FamilyName    string `form:"family_name"`
	Country       string `form:"country"`
	LicenceNumber string `form:"licence_number"`
	Profession    string `form:"profession"`
}
