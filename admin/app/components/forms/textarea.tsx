import React from "react";
import { Label } from "../ui/label";
import { FieldError } from "./field-error";
import type { FieldErrors } from "react-hook-form";
import { Textarea as BaseTextarea } from "~/components/ui/textarea";

interface Props extends React.HTMLProps<HTMLTextAreaElement> {
  name: string;
  label: string;
  errors?: FieldErrors;
  containerClassName?: string;
}

export const Textarea: React.FC<Props> = React.forwardRef(
  ({ errors, name, label, containerClassName, id = name, ...rest }, ref) => {
    const error = errors?.[id]?.message;

    return (
      <div className="flex flex-col gap-1.5">
        <Label htmlFor={id}>{label}</Label>
        <BaseTextarea id={id} name={name} {...rest} ref={ref} />
        <FieldError message={error as string} />
      </div>
    );
  },
);
