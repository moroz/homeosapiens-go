import React from "react";
import { type BlogPost, useListBlogPostsQuery } from "~/hooks";
import { AdminLayout } from "~/components/admin-layout";
import { DataTable } from "~/components/data-table";
import type { ColumnDef } from "@tanstack/react-table";
import { PageTitle } from "~/components/page-title";
import { Link } from "react-router";
import { buttonVariants } from "~/components/ui/button";
import { PlusIcon } from "@phosphor-icons/react";

interface Props {}

const columns: ColumnDef<BlogPost>[] = [
  {
    header: "Title",
    accessorKey: "title",
  },
];

export const BlogPosts: React.FC<Props> = () => {
  const { data: posts, isPending, isError } = useListBlogPostsQuery();

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
        />
      </div>
    </AdminLayout>
  );
};
