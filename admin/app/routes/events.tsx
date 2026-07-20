import { type ColumnDef, type PaginationState } from "@tanstack/react-table";
import { useState } from "react";
import { useNavigate } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { DataTable } from "~/components/data-table";
import { useListEventsQuery } from "~/hooks";
import type { components } from "~/lib/api-types";

type Event = components["schemas"]["Event"];

const dateStyle = { dateStyle: "medium", timeStyle: "short" } as const;

function formatInstant(iso: string) {
  return Temporal.Instant.from(iso).toLocaleString("en-GB", dateStyle);
}

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
      `${formatInstant(row.original.startsAt)} – ${formatInstant(row.original.endsAt)}`,
  },
];

export default function Events() {
  const navigate = useNavigate();
  const [pagination, setPagination] = useState<PaginationState>({ pageIndex: 0, pageSize: 20 });
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
        onPaginationChange={setPagination}
        isPending={isPending}
        isError={isError}
        onRowClick={(event) => navigate(`/events/${event.id}`)}
      />
    </AdminLayout>
  );
}
