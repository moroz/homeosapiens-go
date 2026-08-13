package api

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type userServer struct {
	db *pgxpool.Pool
}

func NewUserServer(db *pgxpool.Pool) *userServer {
	return &userServer{db: db}
}

func (s *userServer) ListUsers(ctx context.Context, params ListUsersRequestObject) (ListUsersResponseObject, error) {
	page, perPage := resolvePaginationParams(params.Params.Page, params.Params.PerPage)

	search := ""
	if params.Params.Search != nil {
		search = *params.Params.Search
	}

	result, err := services.NewUserService(s.db).ListUsers(ctx, &types.ListUsersParams{
		SearchParam: search,
		PerPage:     perPage,
		Page:        page,
	})
	if err != nil {
		return nil, err
	}

	out := make([]User, len(result.Data))
	for i, e := range result.Data {
		out[i] = User{
			Id:               e.ID,
			Email:            openapi_types.Email(e.Email.Plaintext()),
			EmailConfirmedAt: e.EmailConfirmedAt,
			FamilyName:       e.FamilyName.Plaintext(),
			GivenName:        e.GivenName.Plaintext(),
			InsertedAt:       e.InsertedAt,
			PreferredLocale:  string(e.PreferredLocale),
			Role:             UserRole(e.UserRole),
			ProfilePicture:   e.ProfilePicture,
		}
	}

	return ListUsers200JSONResponse{
		Data: out,
		Pagination: Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      result.Pagination.TotalCount,
			TotalPages: result.Pagination.TotalPages,
		},
	}, nil
}

func (s *userServer) GetUser(ctx context.Context, params GetUserRequestObject) (GetUserResponseObject, error) {
	user, err := queries.New(s.db).GetUserByID(ctx, params.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return GetUser404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	return GetUser200JSONResponse{
		Id:               user.ID,
		Email:            openapi_types.Email(user.Email.Plaintext()),
		EmailConfirmedAt: user.EmailConfirmedAt,
		FamilyName:       user.FamilyName.Plaintext(),
		GivenName:        user.GivenName.Plaintext(),
		InsertedAt:       user.InsertedAt,
		PreferredLocale:  string(user.PreferredLocale),
		ProfilePicture:   user.ProfilePicture,
		Role:             UserRole(user.UserRole),
	}, nil
}
