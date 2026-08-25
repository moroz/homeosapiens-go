package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
)

type hostServer struct {
	db queries.DBTX
	q  *queries.Queries
}

func NewHostServer(db queries.DBTX) *hostServer {
	return &hostServer{db: db, q: queries.New(db)}
}

// host maps a host row onto the Host schema, shared by every host operation.
func host(h *queries.Host) Host {
	return Host{
		Id:         h.ID,
		Salutation: h.Salutation,
		GivenName:  h.GivenName,
		FamilyName: h.FamilyName,
		Country:    h.Country,
	}
}

func hostInput(body *HostInput) *types.HostInput {
	return &types.HostInput{
		Salutation: body.Salutation,
		GivenName:  body.GivenName,
		FamilyName: body.FamilyName,
		Country:    body.Country,
	}
}

func (s *hostServer) ListHosts(ctx context.Context, params ListHostsRequestObject) (ListHostsResponseObject, error) {
	page, perPage := resolvePaginationParams(params.Params.Page, params.Params.PerPage)

	hosts, err := s.q.PaginateHosts(ctx, &queries.PaginateHostsParams{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		return nil, err
	}

	count, err := s.q.CountHosts(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]Host, len(hosts))
	for i, h := range hosts {
		out[i] = host(h)
	}

	return ListHosts200JSONResponse{
		Data: out,
		Pagination: Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      count,
			TotalPages: countPages(count, perPage),
		},
	}, nil
}

func (s *hostServer) GetHost(ctx context.Context, request GetHostRequestObject) (GetHostResponseObject, error) {
	h, err := s.q.GetHostById(ctx, request.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return GetHost404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	return GetHost200JSONResponse(host(h)), nil
}

func (s *hostServer) CreateHost(ctx context.Context, request CreateHostRequestObject) (CreateHostResponseObject, error) {
	h, err := services.NewHostService(s.db).CreateHost(ctx, hostInput(request.Body))
	if verr, ok := errors.AsType[validation.Errors](err); ok {
		return CreateHost422JSONResponse{Errors: validationErrorMessages(verr)}, nil
	}
	if err != nil {
		return nil, err
	}

	location := fmt.Sprintf("/api/admin/hosts/%s", h.ID)

	return CreateHost201JSONResponse{
		Headers: CreateHost201ResponseHeaders{
			Location: &location,
		},
		Body: host(h),
	}, nil
}

func (s *hostServer) UpdateHost(ctx context.Context, request UpdateHostRequestObject) (UpdateHostResponseObject, error) {
	h, err := services.NewHostService(s.db).UpdateHost(ctx, request.Id, hostInput(request.Body))
	if errors.Is(err, sql.ErrNoRows) {
		return UpdateHost404Response{}, nil
	}
	if verr, ok := errors.AsType[validation.Errors](err); ok {
		return UpdateHost422JSONResponse{Errors: validationErrorMessages(verr)}, nil
	}
	if err != nil {
		return nil, err
	}

	return UpdateHost200JSONResponse(host(h)), nil
}

// DeleteHost removes a host that nothing points at. Events and videos keep
// their host references, so the foreign keys are what decide whether a host is
// still in use; a violation is reported as a conflict rather than a 500.
func (s *hostServer) DeleteHost(ctx context.Context, request DeleteHostRequestObject) (DeleteHostResponseObject, error) {
	if _, err := s.q.GetHostById(ctx, request.Id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DeleteHost404Response{}, nil
		}
		return nil, err
	}

	err := s.q.DeleteHost(ctx, request.Id)
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23503" {
		return DeleteHost409Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	return DeleteHost204Response{}, nil
}
