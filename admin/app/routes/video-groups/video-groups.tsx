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
    cell: ({ row }) => <span className="font-medium">{row.original.titleEn}</span>,
  },
  { accessorKey: "titlePl", header: "Title (PL)" },
  { accessorKey: "slug", header: "Slug" },
  {
    id: "access",
    header: "Access",
    accessorKey: "isPremium",
    cell: ({ row }) =>
      row.original.isPremium ? (
        <Badge>{`${row.original.price} ${row.original.currency}`}</Badge>
      ) : (
        <Badge variant="secondary">free</Badge>
      ),
  },
  { accessorKey: "videoCount", header: "Videos" },
  {
    id: "insertedAt",
    header: "Created at",
    accessorKey: "insertedAt",
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
    <AdminLayout title="Videos">
      <div className="flex flex-col gap-4">
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold">Video series</h2>
          <Link to="/videos/new" className={buttonVariants()}>
            <PlusIcon />
            New series
          </Link>
        </div>

        <DataTable
          columns={columns}
          data={data?.data ?? []}
          pageCount={data?.pagination.totalPages ?? 0}
          pagination={pagination}
          onPaginationChange={onPaginationChange}
          sorting={sorting}
          onSortingChange={onSortingChange}
          isPending={isPending}
          isError={isError}
          onRowClick={(group) => navigate(`/videos/${group.id}/edit`)}
        />
      </div>
    </AdminLayout>
  );
}
