import React from "react";
import type { FieldErrors } from "react-hook-form";
import { Label } from "~/components/ui/label";
import { Input } from "~/components/ui/input";
import { FieldError } from "~/components/forms/field-error";
import { cn } from "~/lib/utils";

interface Props extends React.HTMLProps<HTMLInputElement> {
  name: string;
  label: string;
  errors?: FieldErrors;
  containerClassName?: string;
}

export const InputField: React.FC<Props> = React.forwardRef(
  ({ errors, name, label, containerClassName, id = name, ...rest }, ref) => {
    const error = errors?.[id]?.message;

    return (
      <div className={cn("flex flex-col gap-1.5", containerClassName)}>
        <Label htmlFor={id}>{label}</Label>
        <Input id={id} name={name} {...rest} ref={ref} />
        <FieldError message={error as string} />
      </div>
    );
  },
);
