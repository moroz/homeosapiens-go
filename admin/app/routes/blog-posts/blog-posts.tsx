import React from "react";
import { type BlogPost, useListBlogPostsQuery } from "~/hooks";
import { AdminLayout } from "~/components/admin-layout";
import { DataTable } from "~/components/data-table";
import type { ColumnDef } from "@tanstack/react-table";
import { PageTitle } from "~/components/page-title";
import { Link, useNavigate } from "react-router";
import { buttonVariants } from "~/components/ui/button";
import { PlusIcon } from "@phosphor-icons/react";
import { formatInstant } from "~/lib/time";
import { Badge } from "~/components/ui/badge";

interface Props {}

const columns: ColumnDef<BlogPost>[] = [
  {
    header: "Title",
    accessorKey: "title",
  },
  {
    header: "Slug",
    accessorKey: "slug",
    cell: ({ row }) => <code>{row.original.slug}</code>,
  },
  {
    header: "Published at",
    accessorKey: "publishedAt",
    cell: ({ row }) => {
      const publishedAt = row.original.publishedAt;
      if (!publishedAt) return <Badge variant="secondary">draft</Badge>;
      return formatInstant(publishedAt);
    },
  },
  {
    header: "Created at",
    accessorKey: "insertedAt",
    cell: ({ row }) => formatInstant(row.original.insertedAt),
  },
];

export const BlogPosts: React.FC<Props> = () => {
  const { data: posts, isPending, isError } = useListBlogPostsQuery();
  const navigate = useNavigate();

  return (
    <AdminLayout title="Blog posts">
      <div className="grid gap-4">
        <header className="flex items-center justify-between">
          <PageTitle>Blog posts</PageTitle>
          <Link to="/blog-posts/new" className={buttonVariants({ variant: "default" })}>
            <PlusIcon />
            New blog post
          </Link>
        </header>

        <DataTable
          columns={columns}
          data={posts ?? []}
          pageCount={1}
          pagination={{ pageIndex: 0, pageSize: 10000 }}
          onPaginationChange={() => null}
          isPending={isPending}
          isError={isError}
          onRowClick={(row) => navigate(row.id)}
        />
      </div>
    </AdminLayout>
  );
};
