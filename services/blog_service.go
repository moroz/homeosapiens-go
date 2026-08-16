package services

import "github.com/moroz/homeosapiens-go/db/queries"

type BlogService struct {
	db queries.DBTX
}

func NewBlogService(db queries.DBTX) *BlogService {
	return &BlogService{db}
}
