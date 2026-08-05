import { PlusIcon } from "@phosphor-icons/react";
import { type ColumnDef } from "@tanstack/react-table";
import { Link, useNavigate } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { DataTable } from "~/components/data-table";
import { buttonVariants } from "~/components/ui/button";
import { useListEventsQuery, useTableSearchParams } from "~/hooks";
import type { components } from "~/lib/api-types";
import { formatInstant } from "~/lib/time";
import { Badge } from "~/components/ui/badge";

type Event = components["schemas"]["Event"];

const columns: ColumnDef<Event>[] = [
  {
    accessorKey: "titleEn",
    header: "Title (EN)",
    cell: ({ row }) => <span className="font-medium">{row.original.titleEn}</span>,
  },
  { accessorKey: "titlePl", header: "Title (PL)" },
  {
    id: "publishedAt",
    header: "Published at",
    accessorKey: "publishedAt",
    cell: ({ row }) => {
      const ts = row.original.publishedAt;
      return ts ? formatInstant(ts) : <Badge variant="secondary">draft</Badge>;
    },
  },
  { accessorKey: "eventType", header: "Type" },
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
      <div className="flex flex-col gap-4">
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold">Events</h2>
          <Link to="/events/new" className={buttonVariants()}>
            <PlusIcon />
            New event
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
          onRowClick={(event) => navigate(`/events/${event.id}`)}
        />
      </div>
    </AdminLayout>
  );
}
