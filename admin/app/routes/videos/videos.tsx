import type { ColumnDef } from "@tanstack/react-table";
import { useNavigate } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { DataTable } from "~/components/data-table";
import { PageTitle } from "~/components/page-title";
import { Badge } from "~/components/ui/badge";
import { useListVideosQuery, useTableSearchParams, type Video } from "~/hooks";

const columns: ColumnDef<Video>[] = [
  {
    accessorKey: "titleEn",
    header: "Title (EN)",
    cell: ({ row }) => (
      <span className="block w-full truncate font-medium text-ellipsis">
        {row.original.titleEn}
      </span>
    ),
  },
  {
    accessorKey: "titlePl",
    header: "Title (PL)",
    cell: ({ row }) => (
      <span className="block w-full truncate text-ellipsis">{row.original.titlePl}</span>
    ),
  },
  {
    accessorKey: "slug",
    header: "Slug",
    cell: ({ row }) => (
      <span className="block w-full truncate font-mono text-sm text-ellipsis">
        {row.original.slug}
      </span>
    ),
  },
  { accessorKey: "provider", header: "Provider", size: 60 },
  {
    accessorKey: "isPublic",
    header: "Visibility",
    size: 60,
    cell: ({ row }) =>
      row.original.isPublic ? (
        <Badge variant="secondary">public</Badge>
      ) : (
        <Badge>members only</Badge>
      ),
  },
  {
    accessorKey: "recordedOn",
    header: "Recorded on",
    size: 70,
    cell: ({ row }) => row.original.recordedOn?.slice(0, 10) ?? null,
  },
];

export default function Videos() {
  const navigate = useNavigate();
  const { pagination, onPaginationChange, sorting, onSortingChange } = useTableSearchParams(20);
  const { data, isPending, isError } = useListVideosQuery(
    pagination.pageIndex + 1,
    pagination.pageSize,
  );

  return (
    <AdminLayout title="Videos">
      <div className="flex flex-col gap-4">
        <header className="flex items-center justify-between">
          <PageTitle>Videos</PageTitle>
        </header>

        <div className="w-full overflow-hidden">
          <DataTable
            className="table-fixed"
            columns={columns}
            data={data?.data ?? []}
            pageCount={data?.pagination.totalPages ?? 0}
            pagination={pagination}
            onPaginationChange={onPaginationChange}
            sorting={sorting}
            onSortingChange={onSortingChange}
            isPending={isPending}
            isError={isError}
            onRowClick={(video) => navigate(`/videos/${video.id}/edit`)}
          />
        </div>
      </div>
    </AdminLayout>
  );
}
