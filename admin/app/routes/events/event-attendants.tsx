import type { ColumnDef } from "@tanstack/react-table";
import React from "react";
import { useParams } from "react-router";
import { AdminLayout } from "~/components/admin-layout";
import { BackButton } from "~/components/back-button";
import { DataTable } from "~/components/data-table";
import { PageTitle } from "~/components/page-title";
import { useGetEventQuery, useTableSearchParams } from "~/hooks";
import { useListEventAttendantsQuery, type EventAttendant } from "~/hooks/event-registrations";

interface Props {}

const columns: ColumnDef<EventAttendant>[] = [
  {
    accessorKey: "givenName",
    header: "Given name",
  },
  {
    accessorKey: "familyName",
    header: "Family name",
  },
  {
    accessorKey: "email",
    header: "Email",
  },
];

export const EventAttendants: React.FC<Props> = () => {
  const { id } = useParams();

  const { data: event, isPending: eventPending } = useGetEventQuery(id!);
  const { onPaginationChange, pagination, sorting, onSortingChange } = useTableSearchParams(20);
  const { data: attendants, isPending: attendantsPending } = useListEventAttendantsQuery({
    eventId: id!,
    page: pagination.pageIndex + 1,
    perPage: pagination.pageSize,
  });

  return (
    <AdminLayout title="Enrolled students">
      <BackButton href={`/events/${id}`}>Back to event</BackButton>
      <div className="grid gap-4">
        {!eventPending && <PageTitle subtitle="Enrolled students">{event?.titleEn}</PageTitle>}

        <DataTable
          columns={columns}
          data={attendants?.data ?? []}
          pageCount={attendants?.pagination.totalPages ?? 0}
          onPaginationChange={onPaginationChange}
          sorting={sorting}
          onSortingChange={onSortingChange}
          pagination={pagination}
          isPending={attendantsPending}
        />
      </div>
    </AdminLayout>
  );
};
