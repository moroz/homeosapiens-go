package types

type Pagination struct {
	Page       int32
	PerPage    int32
	TotalPages int32
	TotalCount int64
}

type PaginationPage[T any] struct {
	Pagination Pagination
	Data       []T
}
