package types

import "github.com/moroz/homeosapiens-go/db/queries"

// CreateUserParams represents the trusted data of a user, coming from a trusted source, such as database data exports. It is intended for backfilling user data, or for seeding the database.
type CreateUserParams struct {
	GivenName  string
	FamilyName string
	Email      string
	Role       queries.UserRole
	Country    string
	Password   string
}
