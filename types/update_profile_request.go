package types

type UpdateProfileRequest struct {
	GivenName  string `form:"given_name"`
	FamilyName string `form:"family_name"`
}
