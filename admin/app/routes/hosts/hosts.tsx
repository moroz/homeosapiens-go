import type { ColumnDef } from "@tanstack/react-table";
import { PlusIcon } from "@phosphor-icons/react";
import { Link, useNavigate } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { DataTable } from "~/components/data-table";
import { PageTitle } from "~/components/page-title";
import { buttonVariants } from "~/components/ui/button";
import { useListHostsQuery, useTableSearchParams, type Host } from "~/hooks";

import { SALUTATION_LABELS } from "./interfaces";

const intl = new Intl.DisplayNames("en", { type: "region" });

const columns: ColumnDef<Host>[] = [
  {
    accessorKey: "familyName",
    header: "Family name",
    cell: ({ row }) => <span className="font-medium">{row.original.familyName}</span>,
  },
  { accessorKey: "givenName", header: "Given name" },
  {
    accessorKey: "salutation",
    header: "Salutation",
    size: 80,
    cell: ({ row }) =>
      row.original.salutation ? SALUTATION_LABELS[row.original.salutation] : null,
  },
  {
    accessorKey: "country",
    header: "Country",
    size: 60,
    cell: ({ row }) => row.original.country && intl.of(row.original.country),
  },
];

export default function Hosts() {
  const navigate = useNavigate();
  const { pagination, onPaginationChange, sorting, onSortingChange } = useTableSearchParams(20);
  const { data, isPending, isError } = useListHostsQuery(
    pagination.pageIndex + 1,
    pagination.pageSize,
  );

  return (
    <AdminLayout title="Hosts">
      <div className="flex flex-col gap-4">
        <header className="flex items-center justify-between">
          <PageTitle>Hosts</PageTitle>
          <Link to="/hosts/new" className={buttonVariants()}>
            <PlusIcon />
            New host
          </Link>
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
            onRowClick={(host) => navigate(`/hosts/${host.id}/edit`)}
          />
        </div>
      </div>
    </AdminLayout>
  );
}
