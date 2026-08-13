import React, { useCallback } from "react";
import { cn } from "~/lib/utils";
import { Button } from "~/components/ui/button";
import { CopyIcon } from "@phosphor-icons/react";

interface DataTableProps {
  children?: React.ReactNode;
  className?: string;
}

export function DetailsTable({ children, className }: DataTableProps) {
  return (
    <table className={cn("details-table", className)}>
      <tbody>{children}</tbody>
    </table>
  );
}

interface DetailsTableFieldProps {
  label: string;
  children: React.ReactNode;
  className?: string;
  copy?: boolean;
}

export function DetailsTableField({ label, children, className, copy }: DetailsTableFieldProps) {
  const onCopy = useCallback(() => {
    navigator.clipboard.writeText(String(children));
  }, [children]);

  return (
    <tr>
      <th className="w-48">{label}</th>
      <td className={className}>
        {children ? (
          <div className="flex items-center gap-3">
            {children}
            {copy && (
              <Button
                variant="ghost"
                type="button"
                size="xs"
                className="font-sans"
                onClick={onCopy}
              >
                <CopyIcon className="w-4" />
                Copy to clipboard
              </Button>
            )}
          </div>
        ) : (
          <span className="text-sm text-muted-foreground">(empty)</span>
        )}
      </td>
    </tr>
  );
}
