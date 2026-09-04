import type { ColumnDef } from "@tanstack/react-table";
import { useNavigate } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { DataTable } from "~/components/data-table";
import { PageTitle } from "~/components/page-title";
import { useListOrdersQuery, useTableSearchParams, type Order } from "~/hooks";
import { formatPrice } from "~/lib/money";
import { formatInstant } from "~/lib/time";

import { OrderStatusBadge } from "./order-status-badge";

const columns: ColumnDef<Order>[] = [
  {
    accessorKey: "orderNumber",
    header: "No.",
    size: 60,
    cell: ({ row }) => <code className="text-sm">{row.original.orderNumber}</code>,
  },
  {
    id: "customer",
    header: "Customer",
    cell: ({ row }) => (
      <div className="flex flex-col">
        <span className="truncate font-medium text-ellipsis">
          {row.original.givenName} {row.original.familyName}
        </span>
        <span className="truncate text-sm text-ellipsis text-muted-foreground">
          {row.original.email}
        </span>
      </div>
    ),
  },
  {
    accessorKey: "billingCountry",
    header: "Country",
    size: 50,
  },
  {
    accessorKey: "status",
    header: "Status",
    size: 60,
    cell: ({ row }) => <OrderStatusBadge status={row.original.status} />,
  },
  {
    accessorKey: "grandTotal",
    header: "Total",
    size: 70,
    cell: ({ row }) => formatPrice(row.original.grandTotal, row.original.currency),
  },
  {
    accessorKey: "insertedAt",
    header: "Placed at",
    size: 80,
    cell: ({ row }) => formatInstant(row.original.insertedAt),
  },
];

export default function Orders() {
  const navigate = useNavigate();
  const { pagination, onPaginationChange, sorting, onSortingChange } = useTableSearchParams(20);
  const { data, isPending, isError } = useListOrdersQuery(
    pagination.pageIndex + 1,
    pagination.pageSize,
  );

  return (
    <AdminLayout title="Orders">
      <div className="flex flex-col gap-4">
        <header className="flex items-center justify-between">
          <PageTitle>Orders</PageTitle>
        </header>

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
            onRowClick={(order) => navigate(`/orders/${order.id}`)}
          />
        </div>
      </div>
    </AdminLayout>
  );
}
