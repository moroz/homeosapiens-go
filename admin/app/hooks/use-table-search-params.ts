import {
  type OnChangeFn,
  type PaginationState,
  type SortingState,
} from "@tanstack/react-table";
import { useCallback } from "react";
import { useSearchParams } from "react-router";

/**
 * Mirrors table pagination and sorting in the URL query string
 * (`?page=&perPage=&sort=&dir=`), so table state survives reloads, deep links,
 * and back/forward. Page is 1-based in the URL, 0-based in TanStack's
 * {@link PaginationState}. Sorting is single-column.
 */
export function useTableSearchParams(defaultPageSize = 20): {
  pagination: PaginationState;
  onPaginationChange: OnChangeFn<PaginationState>;
  sorting: SortingState;
  onSortingChange: OnChangeFn<SortingState>;
} {
  const [searchParams, setSearchParams] = useSearchParams();

  const pageIndex = Math.max(0, (Number(searchParams.get("page")) || 1) - 1);
  const pageSize = Number(searchParams.get("perPage")) || defaultPageSize;
  const pagination: PaginationState = { pageIndex, pageSize };

  const sortId = searchParams.get("sort");
  const sorting: SortingState = sortId
    ? [{ id: sortId, desc: searchParams.get("dir") === "desc" }]
    : [];

  const onPaginationChange = useCallback<OnChangeFn<PaginationState>>(
    (updater) => {
      const next = typeof updater === "function" ? updater(pagination) : updater;
      setSearchParams(
        (prev) => {
          const params = new URLSearchParams(prev);
          params.set("page", String(next.pageIndex + 1));
          params.set("perPage", String(next.pageSize));
          return params;
        },
        { replace: true },
      );
    },
    [pagination, setSearchParams],
  );

  const onSortingChange = useCallback<OnChangeFn<SortingState>>(
    (updater) => {
      const next = typeof updater === "function" ? updater(sorting) : updater;
      setSearchParams(
        (prev) => {
          const params = new URLSearchParams(prev);
          if (next.length > 0) {
            params.set("sort", next[0].id);
            params.set("dir", next[0].desc ? "desc" : "asc");
          } else {
            params.delete("sort");
            params.delete("dir");
          }
          return params;
        },
        { replace: true },
      );
    },
    [sorting, setSearchParams],
  );

  return { pagination, onPaginationChange, sorting, onSortingChange };
}
