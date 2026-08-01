import React from "react";
import { Label } from "../ui/label";
import { FieldError } from "./field-error";
import {
  Controller,
  type Control,
  type FieldErrors,
  type FieldValues,
  type Path,
} from "react-hook-form";
import {
  Select as BaseSelect,
  SelectTrigger,
  SelectValue,
  SelectItem,
  SelectContent,
} from "~/components/ui/select";

interface SelectOption {
  value: string;
  label: string;
}

interface Props<T extends FieldValues> {
  options: SelectOption[];
  label: string;
  name: Path<T>;
  id?: string;
  control: Control<T>;
  errors?: FieldErrors;
}

export function Select<T extends FieldValues>({
  name,
  id = name,
  options,
  label,
  control,
  errors,
}: Props<T>) {
  const error = errors?.[id]?.message as string;

  return (
    <div className="flex flex-col gap-1.5">
      <Label htmlFor={id}>{label}</Label>
      <Controller
        name={name}
        control={control}
        render={({ field }) => (
          <BaseSelect value={field.value}>
            <SelectTrigger id={id} className="w-full">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {options.map((option) => (
                <SelectItem key={option.value} value={option.value}>
                  {option.label}
                </SelectItem>
              ))}
            </SelectContent>
          </BaseSelect>
        )}
      />
      <FieldError message={error} />
    </div>
  );
}
