package types

import (
	"time"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/shopspring/decimal"
)

type VideoGroupListDTO struct {
	*queries.VideoGroup
	HasAccess     bool
	MinRecordedOn *time.Time
	MaxRecordedOn *time.Time
	// Price and Currency are the price of the group's product, so that a locked
	// group can be sold on the spot. Both are nil for a free group.
	Price    *decimal.Decimal
	Currency *string
	// CountInCart is how many times the group's product sits in the visitor's
	// cart, mirroring how paid events render their add-to-cart button.
	CountInCart int
}

type VideoGroupDetailsDTO struct {
	*queries.VideoGroup
	HasAccess   bool
	Videos      []*queries.Video
	Price       *decimal.Decimal
	Currency    *string
	CountInCart int
}

type VideoDetailsDTO struct {
	*queries.Video
	HasAccess bool
	Sources   []*queries.VideoSource
}

type VideoListDTO struct {
	*queries.Video
	Hosts []*queries.Host
}

func (v *VideoGroupListDTO) IsPremium() bool {
	return v.VideoGroup.ProductID != nil
}

func (v *VideoGroupDetailsDTO) IsPremium() bool {
	return v.VideoGroup.ProductID != nil
}
