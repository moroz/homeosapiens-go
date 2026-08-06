package types

type Pagination struct {
	Page       int
	PerPage    int
	TotalPages int
	TotalCount int
}

type PaginationPage[T any] struct {
	Pagination Pagination
	Data       []T
}
