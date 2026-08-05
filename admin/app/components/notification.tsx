import React from "react";
import { Alert, AlertDescription, AlertTitle } from "~/components/ui/alert";
import { InfoIcon, WarningCircleIcon } from "@phosphor-icons/react";
import { cva } from "class-variance-authority";
import { cn } from "~/lib/utils";

interface Props {
  children?: React.ReactNode;
  title: string;
  variant?: "default" | "destructive";
  className?: string;
}

const variants = cva("", {
  variants: {
    variant: {
      destructive:
        "bg-red-50 dark:bg-red-900/45 border-red-200 text-red-950 dark:text-red-50 dark:border-red-100/30",
      default: "",
    },
  },
});

const IconMapping = {
  default: InfoIcon,
  destructive: WarningCircleIcon,
};

export const Notification: React.FC<Props> = ({
  title,
  children,
  variant = "default",
  className,
}) => {
  const Icon = IconMapping[variant] ?? InfoIcon;

  return (
    <Alert className={cn(variants({ variant }), className)}>
      <Icon />
      <AlertTitle>{title}</AlertTitle>
      {children ? <AlertDescription>{children}</AlertDescription> : null}
    </Alert>
  );
};
