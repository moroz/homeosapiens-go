import React from "react";
import { cn } from "~/lib/utils";

interface Props {
  className?: string;
  subtitle?: string;
  children?: React.ReactNode;
}

export const PageTitle: React.FC<Props> = ({ children, className, subtitle }) => {
  return (
    <div className={cn("flex flex-col", className)}>
      <h2 className="text-2xl font-bold">{children}</h2>
      {subtitle ? <p className="subtitle text-xl text-muted-foreground">{subtitle}</p> : null}
    </div>
  );
};
