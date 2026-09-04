import React from "react";
import { Button } from "~/components/ui/button";

interface Props {
  pageIndex: number;
  pageCount: number;
  setPageIndex: (index: number) => void;
  /** Total number of rows on the server, across all pages. */
  total?: number;
  /** Number of rows rendered on the current page. */
  rowCount?: number;
  pageSize?: number;
}

export const Pagination: React.FC<Props> = ({
  pageIndex,
  pageCount,
  setPageIndex,
  total,
  rowCount,
  pageSize,
}) => {
  const firstEntry = pageIndex * (pageSize ?? 0) + 1;
  const lastEntry = firstEntry + (rowCount ?? 0) - 1;

  return (
    <div className="flex items-center justify-between">
      <div className="text-sm text-muted-foreground">
        {total !== undefined && rowCount ? (
          <>
            Entries {firstEntry}&ndash;{lastEntry} of {total}
          </>
        ) : (
          <>
            Page {pageIndex + 1} of {Math.max(pageCount, 1)}
          </>
        )}
      </div>

      {pageCount > 1 && (
        <div className="flex items-center gap-2">
          {Array.from({ length: pageCount }).map((_, i) => {
            return (
              <Button
                variant="outline"
                size="icon"
                className="flex size-8"
                onClick={() => setPageIndex(i)}
              >
                {i + 1}
              </Button>
            );
          })}
        </div>
      )}
    </div>
  );
};
