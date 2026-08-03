import React from "react";
import { Link } from "react-router";
import { buttonVariants } from "~/components/ui/button";
import { CaretLeftIcon as CaretLeft } from "@phosphor-icons/react";

interface Props {
  href: string;
  children: React.ReactNode;
}

export const BackButton: React.FC<Props> = ({ href, children }) => {
  return (
    <Link
      to={href}
      className={buttonVariants({ size: "sm", variant: "ghost", className: "w-fit" })}
    >
      <CaretLeft />
      {children}
    </Link>
  );
};
