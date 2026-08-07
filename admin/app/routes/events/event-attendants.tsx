import { UserPlusIcon } from "@phosphor-icons/react/dist/ssr";
import type { ColumnDef } from "@tanstack/react-table";
import React from "react";
import { Link, Outlet, useParams } from "react-router";
import { AdminLayout } from "~/components/admin-layout";
import { BackButton } from "~/components/back-button";
import { DataTable } from "~/components/data-table";
import { PageTitle } from "~/components/page-title";
import { Button, buttonVariants } from "~/components/ui/button";
import { useGetEventQuery, useTableSearchParams } from "~/hooks";
import { useListEventAttendantsQuery, type EventAttendant } from "~/hooks/event-registrations";
import { formatInstant } from "~/lib/time";

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
  {
    header: "Order #",
    cell: ({ row }) => {
      if (!row.original.orderId) {
        return "Manually enrolled";
      }
      return row.original.orderNumber;
    },
  },
  {
    accessorKey: "insertedAt",
    header: "Enrolled at",
    cell: ({ row }) => formatInstant(row.original.insertedAt),
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
      <div className="grid gap-3">
        <BackButton href={`/events/${id}`}>Back to event</BackButton>
        <div className="grid gap-4">
          <header className="flex justify-between">
            {!eventPending && <PageTitle subtitle="Enrolled students">{event?.titleEn}</PageTitle>}
            <div className="flex items-center">
              <Link to="add" className={buttonVariants({ variant: "outline" })}>
                <UserPlusIcon className="w-5" />
                Add an attendant
              </Link>
            </div>
          </header>

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
      </div>
      <Outlet />
    </AdminLayout>
  );
};
