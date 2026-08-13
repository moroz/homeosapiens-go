import React from "react";
import type { FieldErrors } from "react-hook-form";
import { Label } from "~/components/ui/label";
import { Input } from "~/components/ui/input";
import { Field, FieldLabel, FieldError } from "~/components/ui/field";
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
      <Field className={containerClassName}>
        <FieldLabel htmlFor={id}>{label}</FieldLabel>
        <Input id={id} name={name} {...rest} ref={ref} />
        {error ? <FieldError>{String(error)}</FieldError> : null}
      </Field>
    );
  },
);
