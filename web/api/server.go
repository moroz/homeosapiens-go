package api

import (
	"context"
	"net/http"

	"github.com/moroz/homeosapiens-go/db/queries"
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
// registered under baseURL (e.g. "/admin/api"). The returned handler is wrapped
// into Echo via echo.WrapHandler.
func (s *Server) Handler(baseURL string) http.Handler {
	return HandlerWithOptions(NewStrictHandler(s, nil), StdHTTPServerOptions{
		BaseURL: baseURL,
	})
}

func (s *Server) GetHealth(_ context.Context, _ GetHealthRequestObject) (GetHealthResponseObject, error) {
	return GetHealth200JSONResponse{Status: "ok"}, nil
}

func (s *Server) ListHosts(ctx context.Context, _ ListHostsRequestObject) (ListHostsResponseObject, error) {
	hosts, err := s.q.ListHosts(ctx)
	if err != nil {
		return nil, err
	}

	out := make(ListHosts200JSONResponse, len(hosts))
	for i, h := range hosts {
		out[i] = Host{
			Id:         h.ID,
			GivenName:  h.GivenName,
			FamilyName: h.FamilyName,
			Salutation: h.Salutation,
			Country:    h.Country,
		}
	}
	return out, nil
}
