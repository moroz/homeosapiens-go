import React from "react";
import { cn } from "~/lib/utils";

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

interface DataTableFieldProps {
  label: string;
  children: React.ReactNode;
  className?: string;
}

export function DataTableField({ label, children, className }: DataTableFieldProps) {
  return (
    <tr>
      <th className="w-48">{label}</th>
      <td className={className}>{children}</td>
    </tr>
  );
}
