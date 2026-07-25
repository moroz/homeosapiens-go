import React from "react";

interface DataTableProps {
  children?: React.ReactNode;
}

export function DetailsTable({ children }: DataTableProps) {
  return (
    <table className="data-table">
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
      <th className="w-64">{label}</th>
      <td className={className}>{children}</td>
    </tr>
  );
}
