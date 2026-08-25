package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
)

type HostService struct {
	db queries.DBTX
}

func NewHostService(db queries.DBTX) *HostService {
	return &HostService{db}
}

func (s *HostService) CreateHost(ctx context.Context, params *types.HostInput) (*queries.Host, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return queries.New(s.db).InsertHost(ctx, &queries.InsertHostParams{
		Salutation: params.Salutation,
		GivenName:  params.GivenName,
		FamilyName: params.FamilyName,
		Country:    params.Country,
	})
}

func (s *HostService) UpdateHost(ctx context.Context, id uuid.UUID, params *types.HostInput) (*queries.Host, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return queries.New(s.db).UpdateHost(ctx, &queries.UpdateHostParams{
		ID:         id,
		Salutation: params.Salutation,
		GivenName:  params.GivenName,
		FamilyName: params.FamilyName,
		Country:    params.Country,
	})
}
