import React from "react";
import { Button } from "~/components/ui/button";

interface Props {
  pageIndex: number;
  pageCount: number;
  setPageIndex: (index: number) => void;
}

export const Pagination: React.FC<Props> = ({ pageIndex, pageCount, setPageIndex }) => {
  return (
    <div className="flex items-center justify-between">
      <div className="text-sm text-muted-foreground">
        Page {pageIndex + 1} of {Math.max(pageCount, 1)}
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
