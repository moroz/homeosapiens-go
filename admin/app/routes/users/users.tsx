import React, { useCallback } from "react";
import { useListUsersQuery, type User, useTableSearchParams } from "~/hooks";
import { AdminLayout } from "~/components/admin-layout";
import { DataTable } from "~/components/data-table";
import { type ColumnDef } from "@tanstack/react-table";
import { formatInstant } from "~/lib/time";
import { PageTitle } from "~/components/page-title";
import { useNavigate } from "react-router";

interface Props {}

const columns: ColumnDef<User>[] = [
  {
    id: "fullName",
    header: "Full name",
    accessorFn: (row) => `${row.givenName} ${row.familyName}`,
  },
  {
    header: "Email",
    accessorKey: "email",
  },
  {
    header: "Role",
    accessorKey: "role",
  },
  {
    header: "Account verified at",
    accessorKey: "emailConfirmedAt",
    cell: ({ row }) =>
      row.original.emailConfirmedAt ? formatInstant(row.original.emailConfirmedAt) : "Not verified",
  },
];

export const Users: React.FC<Props> = () => {
  const { pagination, onPaginationChange, sorting, onSortingChange } = useTableSearchParams(20);

  const navigate = useNavigate();

  const { data } = useListUsersQuery({
    page: pagination.pageIndex + 1,
    perPage: pagination.pageSize,
  });

  const onRowClick = useCallback((row: User) => navigate(row.id), [navigate]);

  return (
    <AdminLayout title="Users">
      <div className="grid gap-4">
        <PageTitle>Users</PageTitle>
        <DataTable
          columns={columns}
          data={data?.data ?? []}
          pageCount={data?.pagination.totalPages ?? 1}
          pagination={pagination}
          onPaginationChange={onPaginationChange}
          onRowClick={onRowClick}
        />
      </div>
    </AdminLayout>
  );
};
