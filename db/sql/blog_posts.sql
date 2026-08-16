-- name: ListPublishedBlogPostsByLanguage :many
select * from blog_posts p
where p.published_at is not null
and language = $1 order by id desc;

-- name: ListAllBlogPosts :many
select * from blog_posts p order by id desc;

-- name: GetBlogPostBySlug :one
select * from blog_posts where slug = $1;

-- name: GetBlogPostById :one
select * from blog_posts where id = $1;
