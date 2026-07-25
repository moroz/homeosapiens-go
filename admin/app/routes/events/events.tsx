import { type ColumnDef } from "@tanstack/react-table";
import { useNavigate } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { DataTable } from "~/components/data-table";
import { useListEventsQuery, useTableSearchParams } from "~/hooks";
import type { components } from "~/lib/api-types";
import { formatInstant } from "~/lib/time";

type Event = components["schemas"]["Event"];

const columns: ColumnDef<Event>[] = [
  {
    accessorKey: "titleEn",
    header: "Title (EN)",
    cell: ({ row }) => <span className="font-medium">{row.original.titleEn}</span>,
  },
  { accessorKey: "titlePl", header: "Title (PL)" },
  {
    accessorKey: "eventType",
    header: "Type",
  },
  {
    id: "when",
    header: "When",
    accessorKey: "startsAt",
    cell: ({ row }) =>
      `${formatInstant(row.original.startsAt)}–${formatInstant(row.original.endsAt)}`,
  },
  {
    id: "insertedAt",
    header: "Created at",
    accessorKey: "insertedAt",
    cell: ({ row }) => formatInstant(row.original.insertedAt),
  },
];

export default function Events() {
  const navigate = useNavigate();
  const { pagination, onPaginationChange, sorting, onSortingChange } = useTableSearchParams(20);
  const { data, isPending, isError } = useListEventsQuery(
    pagination.pageIndex + 1,
    pagination.pageSize,
  );

  return (
    <AdminLayout title="Events">
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
        onRowClick={(event) => navigate(`/events/${event.id}`)}
        title="Events"
      />
    </AdminLayout>
  );
}
