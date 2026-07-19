package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/swaggest/swgui/v5emb"
)

// Server implements StrictServerInterface, backed by the sqlc query layer.
// It is mounted in web/router/router.go behind the admin session gate, so
// handlers assume the caller is an authenticated admin and perform no
// per-request authorization.
type Server struct {
	q  *queries.Queries
	db queries.DBTX
}

func NewServer(db queries.DBTX) *Server {
	return &Server{q: queries.New(db), db: db}
}

// Handler builds the net/http handler for the admin API, with every route
// registered under baseURL (e.g. "/api/admin"). Alongside the generated API
// routes it serves the OpenAPI document at <baseURL>/openapi.json and an
// interactive Swagger UI at <baseURL>/docs/ (assets embedded via swgui/v5emb,
// so the page makes no external requests). The returned handler is wrapped into
// Echo via echo.WrapHandler.
func (s *Server) Handler(baseURL string) http.Handler {
	mux := http.NewServeMux()

	// Generated API operations (registered onto our mux via BaseRouter).
	HandlerWithOptions(NewStrictHandler(s, nil), StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: mux,
	})

	// The OpenAPI document, served from the embedded spec so it always matches
	// the generated code. Consumed by the Swagger UI and by frontend codegen.
	mux.HandleFunc("GET "+baseURL+"/openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		spec, err := GetSwagger()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(spec)
	})

	// Swagger UI.
	docsPath := baseURL + "/docs/"
	mux.Handle(docsPath, v5emb.New("Homeo sapiens Admin API", baseURL+"/openapi.json", docsPath))

	return mux
}

func (s *Server) GetHealth(_ context.Context, _ GetHealthRequestObject) (GetHealthResponseObject, error) {
	return GetHealth200JSONResponse{Status: "ok"}, nil
}

func (s *Server) ListHosts(ctx context.Context, params ListHostsRequestObject) (ListHostsResponseObject, error) {
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
		out[i] = Host{
			Id:         h.ID,
			GivenName:  h.GivenName,
			FamilyName: h.FamilyName,
			Salutation: h.Salutation,
			Country:    h.Country,
		}
	}
	result := ListHosts200JSONResponse{
		Data: out,
		Pagination: Pagination{
			Page:       *params.Params.Page,
			PerPage:    perPage,
			Total:      count,
			TotalPages: countPages(count, perPage),
		},
	}

	return result, nil
}

func (s *Server) ListVideos(ctx context.Context, params ListVideosRequestObject) (ListVideosResponseObject, error) {
	page, perPage := resolvePaginationParams(params.Params.Page, params.Params.PerPage)

	videos, err := s.q.PaginateVideos(ctx, &queries.PaginateVideosParams{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		return nil, err
	}

	count, err := s.q.CountVideos(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]Video, len(videos))
	for i, v := range videos {
		out[i] = Video{
			Id:      v.ID,
			TitleEn: v.TitleEn,
			TitlePl: v.TitlePl,
		}
	}

	return ListVideos200JSONResponse{
		Data: out,
		Pagination: Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      count,
			TotalPages: countPages(count, perPage),
		},
	}, nil
}

func (s *Server) ListEvents(ctx context.Context, params ListEventsRequestObject) (ListEventsResponseObject, error) {
	page, perPage := resolvePaginationParams(params.Params.Page, params.Params.PerPage)

	events, err := s.q.PaginateEvents(ctx, &queries.PaginateEventsParams{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		return nil, err
	}

	count, err := s.q.CountEvents(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]Event, len(events))
	for i, e := range events {
		out[i] = Event{
			Id:         e.ID,
			Slug:       e.Slug,
			TitleEn:    e.TitleEn,
			TitlePl:    e.TitlePl,
			SubtitleEn: e.SubtitleEn,
			SubtitlePl: e.SubtitlePl,
			EventType:  string(e.EventType),
			IsVirtual:  e.IsVirtual,
			StartsAt:   e.StartsAt,
			EndsAt:     e.EndsAt,
		}
	}

	return ListEvents200JSONResponse{
		Data: out,
		Pagination: Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      count,
			TotalPages: countPages(count, perPage),
		},
	}, nil
}

func (s *Server) ListUsers(ctx context.Context, params ListUsersRequestObject) (ListUsersResponseObject, error) {
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

	out := make([]User, len(result.Users))
	for i, e := range result.Users {
		out[i] = User{
			Id:               e.ID,
			Email:            openapi_types.Email(e.Email.Plaintext()),
			EmailConfirmedAt: e.EmailConfirmedAt,
			FamilyName:       e.FamilyName.Plaintext(),
			GivenName:        e.GivenName.Plaintext(),
			InsertedAt:       e.InsertedAt,
			PreferredLocale:  string(e.PreferredLocale),
			Role:             UserRole(e.UserRole),
		}
	}

	return ListUsers200JSONResponse{
		Data: out,
		Pagination: Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      result.TotalCount,
			TotalPages: countPages(result.TotalCount, perPage),
		},
	}, nil
}
