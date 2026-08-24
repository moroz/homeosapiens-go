import type { ColumnDef } from "@tanstack/react-table";

import { AdminLayout } from "~/components/admin-layout";
import { DataTable } from "~/components/data-table";
import { PageTitle } from "~/components/page-title";
import { Badge } from "~/components/ui/badge";
import {
  useListProductsQuery,
  useTableSearchParams,
  type Product,
  type ProductType,
} from "~/hooks";
import { formatPrice } from "~/lib/money";
import { formatInstant } from "~/lib/time";

const PRODUCT_TYPE_LABELS: Record<ProductType, string> = {
  event: "Event",
  book: "Book",
  video_group: "Playlist",
};

const columns: ColumnDef<Product>[] = [
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
    accessorKey: "productType",
    header: "Type",
    size: 80,
    cell: ({ row }) => (
      <Badge variant="secondary">{PRODUCT_TYPE_LABELS[row.original.productType]}</Badge>
    ),
  },
  {
    accessorKey: "basePrice",
    header: "Base price",
    size: 80,
    cell: ({ row }) => formatPrice(row.original.basePrice, row.original.basePriceCurrency),
  },
  {
    accessorKey: "insertedAt",
    header: "Created at",
    size: 80,
    cell: ({ row }) => formatInstant(row.original.insertedAt),
  },
];

export default function Products() {
  const { pagination, onPaginationChange, sorting, onSortingChange } = useTableSearchParams(20);
  const { data, isPending, isError } = useListProductsQuery(
    pagination.pageIndex + 1,
    pagination.pageSize,
  );

  return (
    <AdminLayout title="Products">
      <div className="flex flex-col gap-4">
        <header className="flex items-center justify-between">
          <PageTitle>Products</PageTitle>
        </header>

        <div className="w-full overflow-hidden">
          <DataTable
            className="table-fixed"
            columns={columns}
            data={data?.data ?? []}
            pageCount={data?.pagination.totalPages ?? 0}
            pagination={pagination}
            onPaginationChange={onPaginationChange}
            sorting={sorting}
            onSortingChange={onSortingChange}
            isPending={isPending}
            isError={isError}
          />
        </div>
      </div>
    </AdminLayout>
  );
}
