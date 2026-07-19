package api

import (
	"context"
	"encoding/json"
	"math"
	"net/http"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/swaggest/swgui/v5emb"
)

// Server implements StrictServerInterface, backed by the sqlc query layer.
// It is mounted in web/router/router.go behind the admin session gate, so
// handlers assume the caller is an authenticated admin and perform no
// per-request authorization.
type Server struct {
	q *queries.Queries
}

func NewServer(db queries.DBTX) *Server {
	return &Server{q: queries.New(db)}
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
	perPage := int32(config.DefaultPageSize)
	if params.Params.PerPage != nil && *params.Params.PerPage > 0 {
		perPage = int32(*params.Params.PerPage)
	}

	hosts, err := s.q.PaginateHosts(ctx, &queries.PaginateHostsParams{
		Page:    params.Params.Page,
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
			TotalPages: int32(math.Ceil(float64(count) / float64(perPage))),
		},
	}

	return result, nil
}

func (s *Server) ListVideos(ctx context.Context, params ListVideosRequestObject) (ListVideosResponseObject, error) {
	perPage := int32(config.DefaultPageSize)
	if params.Params.PerPage != nil && *params.Params.PerPage > 0 {
		perPage = int32(*params.Params.PerPage)
	}

	videos, err := s.q.PaginateVideos(ctx, &queries.PaginateVideosParams{
		Page:    params.Params.Page,
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
			Page:       *params.Params.Page,
			PerPage:    perPage,
			Total:      count,
			TotalPages: int32(math.Ceil(float64(count) / float64(perPage))),
		},
	}, nil
}
