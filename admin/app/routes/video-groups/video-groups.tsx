import { PlusIcon } from "@phosphor-icons/react";
import { type ColumnDef } from "@tanstack/react-table";
import { Link, useNavigate } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { DataTable } from "~/components/data-table";
import { Badge } from "~/components/ui/badge";
import { buttonVariants } from "~/components/ui/button";
import { useListVideoGroupsQuery, useTableSearchParams, type VideoGroup } from "~/hooks";
import { formatInstant } from "~/lib/time";

const columns: ColumnDef<VideoGroup>[] = [
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
  {
    id: "access",
    header: "Access",
    accessorKey: "isPremium",
    size: 60,
    cell: ({ row }) =>
      row.original.isPremium ? (
        <Badge>{`${row.original.price} ${row.original.currency}`}</Badge>
      ) : (
        <Badge variant="secondary">free</Badge>
      ),
  },
  { accessorKey: "videoCount", header: "Videos", size: 50 },
  {
    id: "insertedAt",
    header: "Created at",
    accessorKey: "insertedAt",
    size: 80,
    cell: ({ row }) => formatInstant(row.original.insertedAt),
  },
];

export default function VideoGroups() {
  const navigate = useNavigate();
  const { pagination, onPaginationChange, sorting, onSortingChange } = useTableSearchParams(20);
  const { data, isPending, isError } = useListVideoGroupsQuery(
    pagination.pageIndex + 1,
    pagination.pageSize,
  );

  return (
    <AdminLayout title="Playlists">
      <div className="flex flex-col gap-4">
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold">Playlists</h2>
          <Link to="/playlists/new" className={buttonVariants()}>
            <PlusIcon />
            New series
          </Link>
        </div>

        <div className="w-full overflow-hidden">
          <DataTable
            className="table-fixed"
            columns={columns}
            data={data?.data ?? []}
            pageCount={data?.pagination.totalPages ?? 0}
            total={data?.pagination.total}
            pagination={pagination}
            onPaginationChange={onPaginationChange}
            sorting={sorting}
            onSortingChange={onSortingChange}
            isPending={isPending}
            isError={isError}
            onRowClick={(group) => navigate(`/playlists/${group.id}/edit`)}
          />
        </div>
      </div>
    </AdminLayout>
  );
}
