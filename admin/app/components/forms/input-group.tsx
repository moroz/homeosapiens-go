import React from "react";
import { cn } from "~/lib/utils";

interface Props {
  columns?: number;
  className?: string;
  children: React.ReactNode;
}

export const InputGroup: React.FC<Props> = ({ columns = 2, children, className }) => {
  return (
    <div
      className={cn("grid gap-4", className)}
      style={{ gridTemplateColumns: `repeat(${columns}, 1fr)` }}
    >
      {children}
    </div>
  );
};
