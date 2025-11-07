package services

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/moroz/homeosapiens-go/db/queries"
)

type EventService struct {
	db queries.DBTX
}

func NewEventService(db queries.DBTX) *EventService {
	return &EventService{db}
}

type EventListDto struct {
	*queries.ListEventsRow
	Hosts []*queries.ListHostsForEventsRow
}

func (s *EventService) ListEvents(ctx context.Context) ([]*EventListDto, error) {
	events, err := queries.New(s.db).ListEvents(ctx)
	if err != nil {
		return nil, err
	}

	var ids []pgtype.UUID
	for _, event := range events {
		ids = append(ids, event.ID)
	}

	hosts, err := queries.New(s.db).ListHostsForEvents(ctx, ids)
	if err != nil {
		return nil, err
	}

	hostMap := make(map[pgtype.UUID][]*queries.ListHostsForEventsRow)
	for _, row := range hosts {
		hostMap[row.EventID] = append(hostMap[row.EventID], row)
	}

	var result []*EventListDto
	for _, event := range events {
		result = append(result, &EventListDto{
			ListEventsRow: event,
			Hosts:         hostMap[event.ID],
		})
	}

	return result, nil
}
