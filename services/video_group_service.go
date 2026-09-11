package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/shopspring/decimal"
)

// VideoGroupService backs the admin API. The public site reads video groups
// through VideoService, which additionally resolves per-user access.
type VideoGroupService struct {
	db queries.DBTX
}

func NewVideoGroupService(db queries.DBTX) *VideoGroupService {
	return &VideoGroupService{db}
}

// VideoGroupAdminDto is a video group as the admin sees it: with its price and
// the number of videos in it, both of which live outside the video_groups row.
type VideoGroupAdminDto struct {
	*queries.VideoGroup
	Price      *decimal.Decimal
	Currency   *string
	VideoCount int
}

func (d *VideoGroupAdminDto) IsPremium() bool {
	return d.ProductID != nil
}

func (s *VideoGroupService) ListVideoGroups(ctx context.Context, page, perPage int32) ([]*VideoGroupAdminDto, int64, error) {
	rows, err := queries.New(s.db).PaginateVideoGroups(ctx, &queries.PaginateVideoGroupsParams{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		return nil, 0, err
	}

	count, err := queries.New(s.db).CountVideoGroups(ctx)
	if err != nil {
		return nil, 0, err
	}

	out := make([]*VideoGroupAdminDto, len(rows))
	for i, row := range rows {
		out[i] = &VideoGroupAdminDto{
			VideoGroup: &row.VideoGroup,
			Price:      row.BasePriceAmount,
			Currency:   row.BasePriceCurrency,
			VideoCount: int(row.VideoCount),
		}
	}

	return out, count, nil
}

func (s *VideoGroupService) GetVideoGroupById(ctx context.Context, id uuid.UUID) (*VideoGroupAdminDto, error) {
	row, err := queries.New(s.db).GetVideoGroupById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &VideoGroupAdminDto{
		VideoGroup: &row.VideoGroup,
		Price:      row.BasePriceAmount,
		Currency:   row.BasePriceCurrency,
		VideoCount: int(row.VideoCount),
	}, nil
}

func (s *VideoGroupService) CreateVideoGroup(ctx context.Context, params *types.CreateVideoGroupInput) (*VideoGroupAdminDto, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.db.(*pgxpool.Pool).Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// A group is only given a product once it carries a non-zero price; the
	// pricing rules are shared with events.
	productId, err := patchProductPricing(ctx, tx, &patchProductPricingParams{
		Price:       optionalFromPtr(params.Price),
		Currency:    optionalFromPtr(params.Currency),
		ProductType: queries.ProductTypeVideoGroup,
		TitlePl:     params.TitlePl,
		TitleEn:     params.TitleEn,
	})
	if err != nil {
		return nil, err
	}

	group, err := queries.New(tx).InsertVideoGroup(ctx, &queries.InsertVideoGroupParams{
		TitleEn:   params.TitleEn,
		TitlePl:   params.TitlePl,
		Slug:      params.Slug,
		ProductID: productId,
	})
	if err != nil {
		if verr := slugTakenError(err); verr != nil {
			return nil, verr
		}
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &VideoGroupAdminDto{VideoGroup: group, Price: params.Price, Currency: params.Currency}, nil
}

func (s *VideoGroupService) UpdateVideoGroup(ctx context.Context, id uuid.UUID, params *types.PatchVideoGroupInput) (*VideoGroupAdminDto, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	group, err := s.GetVideoGroupById(ctx, id)
	if err != nil {
		return nil, err
	}

	if params.IsEmpty() {
		return group, nil
	}

	tx, err := s.db.(*pgxpool.Pool).Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	titlePl := group.TitlePl
	if params.TitlePl.Set {
		titlePl = params.TitlePl.Value
	}

	titleEn := group.TitleEn
	if params.TitleEn.Set {
		titleEn = params.TitleEn.Value
	}

	// Pricing may attach a brand-new product to the group, in which case the
	// video_groups row has to carry the resulting foreign key.
	productId, err := patchProductPricing(ctx, tx, &patchProductPricingParams{
		ProductID:   group.ProductID,
		Price:       params.Price,
		Currency:    params.Currency,
		ProductType: queries.ProductTypeVideoGroup,
		TitlePl:     titlePl,
		TitleEn:     titleEn,
	})
	if err != nil {
		return nil, err
	}

	// The column list is the contract between the payload and the table, so it
	// is spelled out rather than derived from field names at runtime.
	assignments := []struct {
		column string
		value  types.Optional[string]
	}{
		{"title_en", params.TitleEn},
		{"title_pl", params.TitlePl},
		{"slug", params.Slug},
	}

	var query strings.Builder
	query.WriteString("update video_groups set ")

	var queryVars []any
	for _, assignment := range assignments {
		if !assignment.value.Set {
			continue
		}

		queryVars = append(queryVars, assignment.value.Ptr())
		fmt.Fprintf(&query, "%s = $%d, ", assignment.column, len(queryVars))
	}

	if productId != nil {
		queryVars = append(queryVars, productId)
		fmt.Fprintf(&query, "product_id = $%d, ", len(queryVars))
	}

	queryVars = append(queryVars, id)
	fmt.Fprintf(&query, "updated_at = now() where id = $%d", len(queryVars))

	if _, err := tx.Exec(ctx, query.String(), queryVars...); err != nil {
		if verr := slugTakenError(err); verr != nil {
			return nil, verr
		}
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.GetVideoGroupById(ctx, id)
}

func (s *VideoGroupService) ListVideosInVideoGroup(ctx context.Context, groupId uuid.UUID) ([]*queries.Video, error) {
	// The group is read first so that an unknown ID is a 404 rather than an
	// empty list.
	if _, err := queries.New(s.db).GetVideoGroupById(ctx, groupId); err != nil {
		return nil, err
	}

	return queries.New(s.db).ListVideosForVideoGroup(ctx, groupId)
}

// ReplaceVideosInVideoGroup sets the group's videos to exactly videoIds, in that
// order. Positions are unique per group, so incremental edits would have to
// shuffle rows around a unique index; deleting and reinserting inside one
// transaction is simpler, the same way event hosts are replaced.
func (s *VideoGroupService) ReplaceVideosInVideoGroup(ctx context.Context, groupId uuid.UUID, videoIds []uuid.UUID) ([]*queries.Video, error) {
	if err := validateDistinctIds(videoIds); err != nil {
		return nil, err
	}

	if _, err := queries.New(s.db).GetVideoGroupById(ctx, groupId); err != nil {
		return nil, err
	}

	tx, err := s.db.(*pgxpool.Pool).Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := queries.New(tx).DeleteVideoGroupVideos(ctx, groupId); err != nil {
		return nil, err
	}

	for i, videoId := range videoIds {
		err := queries.New(tx).InsertVideoGroupVideo(ctx, &queries.InsertVideoGroupVideoParams{
			VideoID:      videoId,
			VideoGroupID: groupId,
			Position:     int32(i),
		})
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23503" {
			return nil, validation.Errors{
				"videoIds": validation.NewError("exists", "references a video that does not exist"),
			}
		}
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return queries.New(s.db).ListVideosForVideoGroup(ctx, groupId)
}

func validateDistinctIds(ids []uuid.UUID) error {
	seen := make(map[uuid.UUID]struct{}, len(ids))
	for _, id := range ids {
		if _, duplicate := seen[id]; duplicate {
			return validation.Errors{
				"videoIds": validation.NewError("distinct", "must not contain duplicates"),
			}
		}
		seen[id] = struct{}{}
	}

	return nil
}

// slugTakenError translates the unique violation on video_groups.slug into a
// per-field validation error; uniqueness cannot be checked up front without
// racing another writer.
func slugTakenError(err error) error {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if !ok || pgErr.Code != "23505" || !strings.Contains(pgErr.ConstraintName, "slug") {
		return nil
	}

	return validation.Errors{
		"slug": validation.NewError("unique", "has already been taken"),
	}
}

// optionalFromPtr adapts the create payload's pointers to the Optional values
// the shared pricing helper works with: a nil pointer means "not given".
func optionalFromPtr[T any](value *T) types.Optional[T] {
	if value == nil {
		return types.Optional[T]{}
	}
	return types.Some(*value)
}
